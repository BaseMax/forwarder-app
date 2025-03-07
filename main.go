package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const TIME_LIMIT = 100

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

type PortConfig struct {
	Port    int     `json:"port"`
	Gateway string  `json:"gateway"`
	Routes  []Route `json:"routes"`
}

type Config struct {
	Ports []PortConfig `json:"ports"`
}

var config Config

func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   TIME_LIMIT * time.Second,
				KeepAlive: TIME_LIMIT * time.Second,
			}).DialContext,
			MaxIdleConns:          10000,
			IdleConnTimeout:       TIME_LIMIT * time.Second,
			TLSHandshakeTimeout:   TIME_LIMIT * time.Second,
			ResponseHeaderTimeout: TIME_LIMIT * time.Second,
			ExpectContinueTimeout: TIME_LIMIT * time.Second,
		},
		Timeout: TIME_LIMIT * time.Second,
	}
}

var httpClient = NewHTTPClient()

func loadConfig(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	for _, p := range config.Ports {
		if p.Port <= 0 {
			return fmt.Errorf("invalid port number: %d", p.Port)
		}
		if p.Gateway == "" {
			return fmt.Errorf("missing gateway for port %d", p.Port)
		}
		for _, r := range p.Routes {
			if r.Method == "" || r.Path == "" || r.Target == "" {
				return fmt.Errorf("invalid route entry in port %d", p.Port)
			}
		}
	}
	return nil
}

func handleRequest(target string, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, "http://"+target+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	buffer := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := w.Write(buffer[:n]); writeErr != nil {
				return
			}
			flusher.Flush()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error streaming response: %v", err)
			return
		}
	}
}

func startServer(p PortConfig) {
	mux := http.NewServeMux()
	for _, route := range p.Routes {
		target := route.Target
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handleRequest(target, w, r)
		})
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", p.Port),
		Handler:      mux,
		ReadTimeout:  TIME_LIMIT * time.Second,
		WriteTimeout: TIME_LIMIT * time.Second,
		IdleTimeout:  TIME_LIMIT * time.Second,
	}

	log.Printf("Starting server on port %d\n", p.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server on port %d: %v", p.Port, err)
	}
}

func main() {
	if err := loadConfig("config.json"); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	var wg sync.WaitGroup
	for _, p := range config.Ports {
		wg.Add(1)
		go func(portConfig PortConfig) {
			defer wg.Done()
			startServer(portConfig)
		}(p)
	}
	wg.Wait()
}
