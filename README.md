# Service Exporter

🚀 Kubernetes Service Port Forwarding with ngrok

[![Latest Release](https://img.shields.io/github/v/release/Goalt/service-exporter?label=latest)](https://github.com/Goalt/service-exporter/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/Goalt/service-exporter/total)](https://github.com/Goalt/service-exporter/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Goalt/service-exporter)](https://github.com/Goalt/service-exporter/blob/main/go.mod)

A CLI tool that helps you expose Kubernetes services to the internet using ngrok tunnels with interactive configuration.

## Table of Contents

- [Demo](#demo)
- [Installation](#installation) 
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Prerequisites](#prerequisites)
- [Development](#development)

## Demo

![Service Exporter Demo](https://raw.githubusercontent.com/Goalt/service-exporter/refs/heads/copilot/fix-12/assets/demo.gif)

## Installation

### Option 1: Download Prebuilt Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/Goalt/service-exporter/releases/latest):

#### Manual Download

#### Linux (AMD64)
```bash
# Download and install
curl -L -o service-exporter https://github.com/Goalt/service-exporter/releases/latest/download/serviceexporter-linux-amd64
chmod +x service-exporter
sudo mv service-exporter /usr/local/bin/
```

#### Linux (ARM64)
```bash
# Download and install
curl -L -o service-exporter https://github.com/Goalt/service-exporter/releases/latest/download/serviceexporter-linux-arm64
chmod +x service-exporter
sudo mv service-exporter /usr/local/bin/
```

#### macOS (Intel)
```bash
# Download and install
curl -L -o service-exporter https://github.com/Goalt/service-exporter/releases/latest/download/serviceexporter-darwin-amd64
chmod +x service-exporter
sudo mv service-exporter /usr/local/bin/
```

#### macOS (Apple Silicon)
```bash
# Download and install
curl -L -o service-exporter https://github.com/Goalt/service-exporter/releases/latest/download/serviceexporter-darwin-arm64
chmod +x service-exporter
sudo mv service-exporter /usr/local/bin/
```

#### Windows
1. Download [serviceexporter-windows-amd64.exe](https://github.com/Goalt/service-exporter/releases/latest/download/serviceexporter-windows-amd64.exe)
2. Rename to `service-exporter.exe`
3. Add to your PATH or run from the download directory

#### Verify Installation
```bash
service-exporter
# Should start the interactive configuration prompt
# Press Ctrl+C to exit
```

### Option 2: Build from Source

Requires Go 1.25+ and access to the source code:

```bash
git clone https://github.com/Goalt/service-exporter.git
cd service-exporter
make build
```

Or build directly:
```bash
go build -o service-exporter ./cmd
```

## Quick Start

1. **Install** the binary using one of the methods above
2. **Get your ngrok auth token** from [ngrok.com](https://ngrok.com)
3. **Run the application**:
   ```bash
   service-exporter
   ```
4. **Follow the interactive prompts** to configure and expose your service

## Features

- **Interactive Configuration**: Choose between environment variables or manual parameter input
- **Service Discovery**: Automatically lists available Kubernetes services  
- **Port Forwarding**: Creates secure port forwarding to selected services
- **ngrok Integration**: Exposes local ports via ngrok tunnels for external access
- **Graceful Shutdown**: Properly cleans up resources on exit

## Usage

### Running the Application

```bash
service-exporter
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
service-exporter
```

#### Option 2: Provide Parameters Manually

Select this option to be prompted for each parameter:

1. **Ngrok Auth Token**: Enter your ngrok authentication token (input will be masked for security)
2. **Kubeconfig Path**: Enter the path to your kubeconfig file (or press Enter for default)

### Complete Workflow

1. **Configuration**: Choose your preferred configuration method
2. **Service Selection**: Select from a list of available Kubernetes services
3. **Port Selection**: Choose which port of the selected service to forward
4. **Port Forwarding**: The tool forwards the selected service port to a local port
5. **ngrok Tunnel**: Creates a public URL for external access
6. **Access**: Use the provided public URL to access your service

Example output:
```
🎉 Setup complete!
==================
Service: my-service (ns: default)
Selected Port: 8080 (http)
Local Port: 8080
Public URL: https://abc123.ngrok.io

You can now access your service via the public URL above!

📌 Press Ctrl+C to gracefully shutdown and cleanup resources...
```

## Prerequisites

- Go 1.25+ (for building from source)
- Access to a Kubernetes cluster
- ngrok account and authentication token
- Valid kubeconfig file

## Getting ngrok Auth Token

1. Sign up at [ngrok.com](https://ngrok.com)
2. Go to your [ngrok dashboard](https://dashboard.ngrok.com/get-started/your-authtoken)
3. Copy your authentication token

## Development
Repository contains configs for devcontainer with all necessary setup.

### Building

```bash
make build
# or
go build -o service-exporter ./cmd
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

## Troubleshooting

### Common Issues

**Q: "No services found" error**
```bash
# Check your kubernetes context
kubectl config current-context
kubectl get services --all-namespaces
```

**Q: ngrok authentication fails**
```bash
# Verify your token is set correctly
echo $NGROK_AUTH_TOKEN
# Or test ngrok directly
ngrok config check
```

**Q: Port forwarding fails**
```bash
# Check if service exists and has ports
kubectl describe service <service-name>
# Verify cluster connectivity
kubectl cluster-info
```

**Q: Permission denied installing binary**
```bash
# Install to user directory instead
curl -L -o ~/bin/service-exporter https://github.com/Goalt/service-exporter/releases/latest/download/serviceexporter-linux-amd64
chmod +x ~/bin/service-exporter
export PATH="$HOME/bin:$PATH"
```

## Contributing

We welcome contributions! Here's how you can help:

1. **🐛 Report Issues**: Found a bug? [Open an issue](https://github.com/Goalt/service-exporter/issues/new)
2. **💡 Suggest Features**: Have an idea? [Create a feature request](https://github.com/Goalt/service-exporter/issues/new)
3. **📝 Improve Documentation**: Help make the docs better
4. **🔧 Submit Code**: Fork, develop, and submit a pull request

### Development Setup

```bash
# Clone and setup
git clone https://github.com/Goalt/service-exporter.git
cd service-exporter
make deps

# Run tests
make test

# Build
make build
```
