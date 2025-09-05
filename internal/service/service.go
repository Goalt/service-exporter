package service

import (
	"fmt"
	"math/rand"
	"time"
)

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

// MockService implements the Service interface with mocked functionality
type MockService struct {
	activeService  string
	activePort     int
	activeNgrokURL string
}

// NewMockService creates a new mock service instance
func NewMockService() Service {
	return &MockService{}
}

// GetServices returns a mocked list of Kubernetes services
func (m *MockService) GetServices() ([]string, error) {
	// Simulate some delay
	time.Sleep(500 * time.Millisecond)

	services := []string{
		"web-frontend",
		"api-gateway",
		"user-service",
		"database-service",
		"cache-service",
		"notification-service",
	}

	return services, nil
}

// StartPortForwarding simulates starting port forwarding for a service
func (m *MockService) StartPortForwarding(serviceName string) (int, error) {
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
func (m *MockService) CreateNgrokSession(port int) (string, error) {
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
func (m *MockService) Cleanup() error {
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
