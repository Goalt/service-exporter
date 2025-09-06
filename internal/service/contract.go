package service

import "context"

// Service defines the interface for Kubernetes service operations
type Service interface {
	// GetServices returns a list of available Kubernetes services
	GetServices(ctx context.Context) ([]string, error)

	// StartPortForwarding starts port forwarding for the specified service
	StartPortForwarding(ctx context.Context, serviceName string) (int, error)

	// CreateNgrokSession creates an ngrok session for the forwarded port
	CreateNgrokSession(ctx context.Context, port int) (string, error)

	// Cleanup performs graceful shutdown of all active sessions
	Cleanup() error
}

type K8s interface {
	// ListServices lists all services in the Kubernetes cluster
	ListServices(ctx context.Context) ([]string, error)

	// PortForward creates a port-forward connection to a service
	PortForward(ctx context.Context, serviceName string, namespace string, localPort int) error
}

// NgrokClient defines the interface for ngrok client operations
type NgrokClient interface {
	StartTunnel(ctx context.Context, port int) (string, error)
	Close() error
}
