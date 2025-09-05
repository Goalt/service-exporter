package service

import (
	"fmt"
	"math/rand"
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

// StartPortForwarding simulates starting port forwarding for a service
func (m *service) StartPortForwarding(serviceName string) (int, error) {
	// Simulate some delay
	time.Sleep(1 * time.Second)

	// Generate a random port between 8000-9000
	port := 8000 + rand.Intn(1000)

	// Store the active session info
	m.activeService = serviceName
	m.activePort = port

	fmt.Printf("🔄 Starting port forwarding for service '%s' on port %d...\n", serviceName, port)

	return port, nil
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
