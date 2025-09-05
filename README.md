# Service Exporter - Kubernetes Service Port Forwarding with ngrok

A CLI tool that provides easy port forwarding for Kubernetes services with optional ngrok tunnel creation.

## Features

- **Real Kubernetes Integration**: Connect to live Kubernetes clusters
- **Automatic Service Discovery**: List available services from your cluster
- **Port Forwarding**: Create real port-forwarding connections to Kubernetes pods
- **Fallback Mode**: Mock service for demonstration when cluster is unavailable
- **Graceful Cleanup**: Proper cleanup of port-forwarding sessions on exit

## Usage

### Prerequisites

- Go 1.24+ 
- Access to a Kubernetes cluster (optional, will fallback to mock mode)
- `kubectl` configured with valid kubeconfig (for real cluster access)

### Running the Application

#### With Real Kubernetes Cluster

```bash
# Build and run with default namespace
go run ./cmd/main.go

# Or build and run
go build -o service-exporter ./cmd/main.go
./service-exporter
```

#### With Mock Service (Demo Mode)

```bash
# Explicitly use mock service
USE_MOCK=true go run ./cmd/main.go
```

#### Custom Namespace

```bash
# Use a specific namespace
K8S_NAMESPACE=my-namespace go run ./cmd/main.go
```

### Environment Variables

- `USE_MOCK=true`: Force use of mock service for demonstration
- `K8S_NAMESPACE`: Specify Kubernetes namespace (default: "default")
- `KUBECONFIG`: Path to kubeconfig file (default: ~/.kube/config)

### Kubernetes Configuration

The application will attempt to connect to Kubernetes in this order:

1. **In-cluster config**: When running inside a Kubernetes pod
2. **Kubeconfig file**: From `$KUBECONFIG` or `~/.kube/config`
3. **Fallback to mock**: If no valid configuration is found

## How It Works

1. **Service Discovery**: Lists all services in the specified namespace
2. **Service Selection**: Interactive prompt to choose a service
3. **Pod Discovery**: Finds running pods that match the service selector
4. **Port Forwarding**: Creates a direct port-forward to the pod
5. **Ngrok Tunnel**: Creates a public URL (currently simulated)
6. **Graceful Cleanup**: Stops all connections on exit (Ctrl+C)

## Development

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

### Project Structure

```
├── cmd/main.go                 # Main application entry point
├── internal/
│   ├── service/
│   │   ├── service.go          # Service interface and implementations
│   │   └── service_test.go     # Service tests
│   └── prompt/
│       ├── prompt.go           # Interactive prompts
│       └── prompt_test.go      # Prompt tests
├── go.mod                      # Go module dependencies
└── README.md                   # This file
```

## Implementation Details

### Service Interface

The application uses a `Service` interface that supports both real and mock implementations:

```go
type Service interface {
    GetServices() ([]string, error)
    StartPortForwarding(serviceName string) (int, error)
    CreateNgrokSession(port int) (string, error)
    Cleanup() error
}
```

### Kubernetes Integration

- Uses `k8s.io/client-go` for real Kubernetes API access
- Implements actual port-forwarding using the Kubernetes port-forward API
- Automatically discovers services and their associated pods
- Handles pod selection and port mapping

### Error Handling

- Graceful fallback when Kubernetes cluster is not available
- Proper error messages for common issues (missing kubeconfig, no services, etc.)
- Cleanup on interruption signals (SIGINT, SIGTERM)

## Troubleshooting

### "Failed to connect to Kubernetes cluster"

This is expected when:
- No kubeconfig is available
- Cluster is not accessible
- Invalid credentials

The application will automatically fall back to mock mode for demonstration.

### "No services available"

- Check if you're in the correct namespace
- Verify services exist: `kubectl get services -n <namespace>`
- Ensure your kubeconfig has proper permissions

### Port forwarding fails

- Check if pods are running: `kubectl get pods -n <namespace>`
- Verify service selectors match running pods
- Ensure the service has defined ports