# Service Exporter Examples

This document provides usage examples for the Service Exporter application.

## Example 1: Running with Mock Service (Demo Mode)

When no Kubernetes cluster is available or when you want to test the application:

```bash
USE_MOCK=true ./service-exporter
```

Output:
```
🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok
================================================================
📋 Using mock Kubernetes service for demonstration...

📋 Fetching available Kubernetes services...
? Select a Kubernetes service: 
  ▸ web-frontend
    api-gateway
    user-service
    database-service
    cache-service
    notification-service

✅ Selected service: web-frontend

🔄 Starting port forwarding for service 'web-frontend' on port 8123...

🌐 Creating ngrok tunnel for port 8123...

🎉 Setup complete!
==================
Service: web-frontend
Local Port: 8123
Public URL: https://a1b2c3d4.ngrok.io

You can now access your service via the public URL above!

📌 Press Ctrl+C to gracefully shutdown and cleanup resources...
```

## Example 2: Running with Real Kubernetes Cluster

When connected to a Kubernetes cluster:

```bash
# Use default namespace
./service-exporter

# Or specify a namespace
K8S_NAMESPACE=production ./service-exporter
```

Expected workflow:
1. Application connects to Kubernetes cluster
2. Lists real services from the specified namespace
3. User selects a service
4. Application finds running pods for the service
5. Creates actual port-forwarding to the pod
6. Sets up ngrok tunnel (simulated)

## Example 3: Automatic Fallback

When Kubernetes is not available, the application gracefully falls back:

```bash
./service-exporter
```

Output:
```
🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok
================================================================
📋 Connecting to Kubernetes cluster...
❌ Failed to connect to Kubernetes cluster: failed to get kubernetes config: stat /home/user/.kube/config: no such file or directory
💡 Falling back to mock service for demonstration...
   Set USE_MOCK=true to explicitly use mock service

📋 Fetching available Kubernetes services...
[continues with mock data...]
```

## Example 4: Environment Variable Configuration

```bash
# Set custom kubeconfig path
KUBECONFIG=/path/to/custom/kubeconfig ./service-exporter

# Use specific namespace in production
K8S_NAMESPACE=production ./service-exporter

# Force mock mode for testing
USE_MOCK=true ./service-exporter
```

## Example 5: Building and Testing

```bash
# Build the application
go build -o service-exporter ./cmd/main.go

# Run tests
go test ./...

# Run specific service tests
go test ./internal/service -v

# Run with coverage
go test -cover ./...
```

## Real-World Usage Scenarios

### Development Environment
```bash
# Connect to local minikube cluster
kubectl config use-context minikube
./service-exporter
```

### Production Troubleshooting
```bash
# Access production services for debugging
K8S_NAMESPACE=production ./service-exporter
# Select the problematic service
# Access it locally for debugging
```

### CI/CD Integration
```bash
# In CI pipelines, use mock mode for testing
USE_MOCK=true ./service-exporter &
# Run automated tests against the mock endpoints
```

## Troubleshooting Examples

### Missing Kubeconfig
```bash
$ ./service-exporter
📋 Connecting to Kubernetes cluster...
❌ Failed to connect to Kubernetes cluster: failed to get kubernetes config: stat /home/user/.kube/config: no such file or directory
💡 Falling back to mock service for demonstration...
```

**Solution**: Configure kubectl or set `USE_MOCK=true`

### No Services Found
```bash
$ K8S_NAMESPACE=empty-namespace ./service-exporter
📋 Connecting to Kubernetes cluster...
📋 Fetching available Kubernetes services...
Error: no services available
```

**Solution**: Check namespace and verify services exist:
```bash
kubectl get services -n empty-namespace
```

### Port Forwarding Failure
If port forwarding fails, the application will show specific error messages about pod selection, port availability, or connectivity issues.

## Advanced Configuration

### Custom Kubeconfig
```bash
export KUBECONFIG=/path/to/special/kubeconfig
./service-exporter
```

### Multiple Cluster Support
```bash
# Switch context and run
kubectl config use-context cluster-1
./service-exporter

# Switch to another cluster
kubectl config use-context cluster-2  
K8S_NAMESPACE=staging ./service-exporter
```