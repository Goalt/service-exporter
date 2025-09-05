package service

// Service defines the interface for Kubernetes service operations
type Service interface {
	// GetServices returns a list of available Kubernetes services
	GetServices() ([]string, error)

	// StartPortForwarding starts port forwarding for the specified service
	StartPortForwarding(serviceName string) (int, error)

	// CreateNgrokSession creates an ngrok session for the forwarded port
	CreateNgrokSession(port int) (string, error)

	// Cleanup performs graceful shutdown of all active sessions
	Cleanup() error
}

type K8s interface {
	// ListServices lists all services in the Kubernetes cluster
	ListServices() ([]string, error)
	
	// PortForward creates a port-forward connection to a service
	PortForward(serviceName string, namespace string, localPort int) error
}
