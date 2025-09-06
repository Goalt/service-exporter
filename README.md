# Service Exporter

🚀 Kubernetes Service Port Forwarding with ngrok

A CLI tool that helps you expose Kubernetes services to the internet using ngrok tunnels with interactive configuration.

## Features

- **Interactive Configuration**: Choose between environment variables or manual parameter input
- **Service Discovery**: Automatically lists available Kubernetes services  
- **Port Forwarding**: Creates secure port forwarding to selected services
- **ngrok Integration**: Exposes local ports via ngrok tunnels for external access
- **Graceful Shutdown**: Properly cleans up resources on exit

## Usage

### Running the Application

```bash
go run ./cmd/main.go
```

When you start the application, you'll be presented with a configuration choice:

```
⚙️  Configuration Setup
=====================
? Configuration mode: 
  ▸ Use default values from environment variables
    Provide parameters manually
```

### Configuration Options

#### Option 1: Use Default Values (Environment Variables)

Select this option to use environment variables for configuration:

- `NGROK_AUTH_TOKEN` (required): Your ngrok authentication token
- `KUBECONFIG` (optional): Path to your kubeconfig file (defaults to `~/.kube/config`)

Example:
```bash
export NGROK_AUTH_TOKEN="your_ngrok_token_here"
export KUBECONFIG="/path/to/your/kubeconfig"
go run ./cmd/main.go
```

#### Option 2: Provide Parameters Manually

Select this option to be prompted for each parameter:

1. **Ngrok Auth Token**: Enter your ngrok authentication token (input will be masked for security)
2. **Kubeconfig Path**: Enter the path to your kubeconfig file (or press Enter for default)

### Complete Workflow

1. **Configuration**: Choose your preferred configuration method
2. **Service Selection**: Select from a list of available Kubernetes services
3. **Port Forwarding**: The tool automatically forwards the service to a local port
4. **ngrok Tunnel**: Creates a public URL for external access
5. **Access**: Use the provided public URL to access your service

Example output:
```
🎉 Setup complete!
==================
Service: my-service (ns: default)
Local Port: 8080
Public URL: https://abc123.ngrok.io

You can now access your service via the public URL above!

📌 Press Ctrl+C to gracefully shutdown and cleanup resources...
```

## Prerequisites

- Go 1.24+
- Access to a Kubernetes cluster
- ngrok account and authentication token
- Valid kubeconfig file

## Getting ngrok Auth Token

1. Sign up at [ngrok.com](https://ngrok.com)
2. Go to your [ngrok dashboard](https://dashboard.ngrok.com/get-started/your-authtoken)
3. Copy your authentication token

## Development

### Building

```bash
go build -o serviceexporter ./cmd/main.go
```

### Testing

```bash
go test ./...
```

### Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── k8s/                 # Kubernetes client
│   ├── ngrok/               # ngrok client  
│   ├── prompt/              # Interactive prompts
│   └── service/             # Core service logic
├── go.mod
└── go.sum
```

## Error Handling

The application provides clear error messages for common issues:

- **Missing ngrok token**: When using default configuration without `NGROK_AUTH_TOKEN` set
- **Invalid kubeconfig**: When the specified kubeconfig file is not found or invalid
- **No services found**: When no Kubernetes services are available in the cluster
- **Port forwarding failures**: When unable to establish port forwarding to the selected service
- **ngrok connection issues**: When unable to create ngrok tunnel

## License

[Add your license information here]