# Proxymity

A lightweight, high-performance reverse proxy and load balancer written in Go.

## Features

- ✅ **Reverse Proxy** - Forward HTTP requests to backend servers
- ✅ **Load Balancing** - Distribute traffic using round-robin algorithm
- ✅ **Health Checks** - Automatic backend health monitoring
- ✅ **Configuration** - Simple YAML-based configuration
- ✅ **Graceful Shutdown** - Clean server shutdown on SIGINT/SIGTERM
- ✅ **Error Handling** - Automatic failover when backends are unavailable

## Quick Start

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd proxymity

# Install dependencies
go mod download
```

### Configuration

Create a `config.yaml` file:

```yaml
proxy:
  host: "0.0.0.0"
  port: "8080"

backend:
  - name: "backend-1"
    url: "http://localhost:8081"
    enabled: true
  
  - name: "backend-2"
    url: "http://localhost:8082"
    enabled: true
  
  - name: "backend-3"
    url: "http://localhost:8083"
    enabled: true
```

### Running

```bash
# Start the proxy
go run cmd/main.go

# Or build and run
go build -o proxymity cmd/main.go
./proxymity
```

### Testing

```bash
# Start some dummy backend servers (in separate terminals)
cd dummies
go run cmd/main.go  # Will start on port 8081, 8082, etc.

# Send requests through the proxy
curl http://localhost:8080/
curl http://localhost:8080/health
```

## Project Structure

```
proxymity/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── backend/             # Backend management
│   │   ├── backend.go       # Backend struct and methods
│   │   └── pool.go          # Backend pool management
│   ├── config/              # Configuration handling
│   │   └── config.go        # Config loading and parsing
│   ├── loadbalancer/        # Load balancing strategies
│   │   ├── loadbalancer.go  # LoadBalancer interface
│   │   └── roundrobin.go    # Round-robin implementation
│   ├── proxy/               # Reverse proxy logic
│   │   └── proxy.go         # HTTP request forwarding
│   └── server/              # HTTP server setup
│       └── server.go        # Server initialization
├── dummies/                 # Test backend servers
│   └── cmd/main.go
├── config.yaml              # Configuration file
├── go.mod
└── README.md
```

## Architecture

```
Client Request
      ↓
[Proxymity Reverse Proxy]
      ↓
[Load Balancer (Round-Robin)]
      ↓
[Backend Pool]
   ↓  ↓  ↓
[B1][B2][B3] Backend Servers
```

## Load Balancing

Proxymity uses a **round-robin** load balancing algorithm that:
- Distributes requests evenly across all healthy backends
- Automatically skips unhealthy backends
- Uses atomic operations for thread-safe concurrent access

## Error Handling

- **Backend Unavailable**: Returns 502 Bad Gateway and marks backend as unhealthy
- **No Healthy Backends**: Returns 503 Service Unavailable
- **Automatic Recovery**: Backends can be marked healthy again via health checks (coming soon)

## Roadmap

- [ ] Multiple load balancing strategies (least connections, weighted, IP hash)
- [ ] Active health checks
- [ ] Rate limiting
- [ ] Metrics and monitoring (Prometheus)
- [ ] TLS/SSL termination
- [ ] Response caching
- [ ] Request/response logging middleware
- [ ] WebSocket support
- [ ] Kubernetes integration
- [ ] Hot-reload configuration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - feel free to use this project however you like!

## Author

Built with ❤️ using Go and Gin
