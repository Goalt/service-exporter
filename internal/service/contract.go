package service

import "context"

// ServicePort represents a service port with its details
type ServicePort struct {
	Name       string
	Port       int32
	TargetPort int32
	Protocol   string
}

// Service defines the interface for Kubernetes service operations
type Service interface {
	// GetServices returns a list of available Kubernetes services
	GetServices(ctx context.Context) ([]string, error)

	// GetServicePorts returns available ports for a specific service
	GetServicePorts(ctx context.Context, serviceName string) ([]ServicePort, error)

	// StartPortForwarding starts port forwarding for the specified service and port
	StartPortForwarding(ctx context.Context, serviceName string, servicePort int32) (int, error)

	// CreateNgrokSession creates an ngrok session for the forwarded port
	CreateNgrokSession(ctx context.Context, port int) (string, error)

	// Cleanup performs graceful shutdown of all active sessions
	Cleanup() error
}

type K8s interface {
	// ListServices lists all services in the Kubernetes cluster
	ListServices(ctx context.Context) ([]string, error)

	// GetServicePorts returns available ports for a specific service
	GetServicePorts(ctx context.Context, serviceName string, namespace string) ([]ServicePort, error)

	// PortForward creates a port-forward connection to a service
	PortForward(ctx context.Context, serviceName string, namespace string, localPort int, servicePort int32) error
}

// NgrokClient defines the interface for ngrok client operations
type NgrokClient interface {
	StartTunnel(ctx context.Context, port int) (string, error)
	Close() error
}
