# Forwarder App

Forwarder App is a simple HTTP proxy server written in Go. It forwards requests from different ports to specified target URLs based on predefined routes. This is useful for routing traffic between microservices, load balancing, or general request forwarding.

## Features

- **Multiple Ports**: Supports multiple ports with configurable gateways.
- **Custom Routes**: Forward requests based on HTTP method and path.
- **TLS Support**: Configurable TLS settings for secure communication.
- **Streaming Support**: Handles large responses efficiently with streaming.
- **Concurrency**: Handles multiple ports and routes concurrently.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/BaseMax/forwarder-app.git
   cd forwarder-app
   ```

2. Install Go (if you don't have it):
   [Download Go](https://golang.org/dl/)

3. Run the application:
   ```bash
   go run main.go
   ```

4. Optionally, build the application:
   ```bash
   go build -o forwarder-app.exe
   ```

## Configuration

The application requires a `config.json` file to define the ports and routes. Here is an example of the `config.json` file:

```json
{
  "ports": [
    {
      "port": 9004,
      "gateway": "127.0.0.1:2004",
      "routes": [
        { "method": "POST", "path": "/v1/member/register", "target": "127.0.0.1:30031" },
        { "method": "POST", "path": "/v1/login", "target": "127.0.0.1:30031" },
        { "method": "POST", "path": "/v1/verifycode", "target": "127.0.0.1:30031" },
        { "method": "POST", "path": "/v1/getforgetcode", "target": "127.0.0.1:30031" },
        { "method": "POST", "path": "/v1/newpassword", "target": "127.0.0.1:30031" },
        { "method": "POST", "path": "/v1/resendactive", "target": "127.0.0.1:30031" },
        { "method": "POST", "path": "/v1/activeuser", "target": "127.0.0.1:30031" }
      ]
    },
    {
      "port": 9005,
      "gateway": "127.0.0.1:3005",
      "routes": []
    },
    {
      "port": 9524,
      "gateway": "127.0.0.1:2524",
      "routes": []
    }
  ]
}
```

### Configuration Details:
- `ports`: List of ports the server will listen on.
- `port`: Port number for the server.
- `gateway`: The address of the target server to forward requests to.
- `routes`: List of routes specifying the HTTP method, path, and target server.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Copyright

© 2025 Max Base.
