package service

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// service implements the Service interface for Kubernetes service operations
type service struct {
	activeService  string
	activePort     int
	activeNgrokURL string

	client K8s
}

// NewService creates a new service instance
func NewService(client K8s) Service {
	return &service{
		client: client,
	}
}

// GetServices returns a list of Kubernetes services from the cluster
func (m *service) GetServices() ([]string, error) {
	if m.client == nil {
		return nil, fmt.Errorf("kubernetes client not available")
	}

	return m.client.ListServices()
}

// StartPortForwarding starts real port forwarding for a service
func (m *service) StartPortForwarding(serviceName string) (int, error) {
	if m.client == nil {
		return 0, fmt.Errorf("kubernetes client not available")
	}

	// Parse service name to extract service name and namespace
	// Format: "service-name (ns: namespace)"
	actualServiceName, namespace, err := m.parseServiceName(serviceName)
	if err != nil {
		return 0, fmt.Errorf("failed to parse service name: %w", err)
	}

	// Find an available local port
	localPort, err := m.findAvailablePort()
	if err != nil {
		return 0, fmt.Errorf("failed to find available port: %w", err)
	}

	// Start port forwarding using the Kubernetes client
	fmt.Printf("🔄 Starting port forwarding for service '%s' in namespace '%s' on local port %d...\n", actualServiceName, namespace, localPort)
	
	err = m.client.PortForward(actualServiceName, namespace, localPort)
	if err != nil {
		return 0, fmt.Errorf("failed to start port forwarding: %w", err)
	}

	// Store the active session info
	m.activeService = serviceName
	m.activePort = localPort

	return localPort, nil
}

// parseServiceName extracts service name and namespace from the formatted string
// Input format: "service-name (ns: namespace)"
func (m *service) parseServiceName(serviceName string) (string, string, error) {
	// Find the namespace part
	nsIndex := strings.Index(serviceName, " (ns: ")
	if nsIndex == -1 {
		// No namespace specified, assume default namespace
		return serviceName, "default", nil
	}

	actualServiceName := serviceName[:nsIndex]
	namespaceStart := nsIndex + 6 // len(" (ns: ")
	namespaceEnd := strings.LastIndex(serviceName, ")")
	
	if namespaceEnd == -1 || namespaceEnd <= namespaceStart {
		return "", "", fmt.Errorf("invalid service name format: %s", serviceName)
	}

	namespace := serviceName[namespaceStart:namespaceEnd]
	return actualServiceName, namespace, nil
}

// findAvailablePort finds an available local port in the range 8000-9000
func (m *service) findAvailablePort() (int, error) {
	for port := 8000; port <= 9000; port++ {
		if m.isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports in range 8000-9000")
}

// isPortAvailable checks if a port is available for use
func (m *service) isPortAvailable(port int) bool {
	address := net.JoinHostPort("localhost", strconv.Itoa(port))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

// CreateNgrokSession simulates creating an ngrok session
func (m *service) CreateNgrokSession(port int) (string, error) {
	// Simulate some delay
	time.Sleep(2 * time.Second)

	// Generate a mock ngrok URL
	randomId := fmt.Sprintf("%x", rand.Uint32())
	ngrokURL := fmt.Sprintf("https://%s.ngrok.io", randomId)

	// Store the active ngrok URL
	m.activeNgrokURL = ngrokURL

	fmt.Printf("🌐 Creating ngrok tunnel for port %d...\n", port)

	return ngrokURL, nil
}

// Cleanup performs graceful shutdown of all active sessions
func (m *service) Cleanup() error {
	fmt.Println("\n🔄 Performing graceful shutdown...")

	if m.activeNgrokURL != "" {
		fmt.Printf("🔌 Closing ngrok tunnel: %s\n", m.activeNgrokURL)
		time.Sleep(500 * time.Millisecond) // Simulate cleanup delay
		m.activeNgrokURL = ""
	}
	if m.activeService != "" && m.activePort > 0 {
		fmt.Printf("🔌 Stopping port forwarding for service '%s' on port %d\n", m.activeService, m.activePort)
		time.Sleep(500 * time.Millisecond) // Simulate cleanup delay
		m.activeService = ""
		m.activePort = 0
	}

	fmt.Println("✅ Graceful shutdown completed")
	return nil
}
