package service

import (
	"context"
	"fmt"
	"testing"
)

// mockK8sClient implements the K8s interface for testing
type mockK8sClient struct {
	services []string
	err      error
}

func (m *mockK8sClient) ListServices(ctx context.Context) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.services, nil
}

func (m *mockK8sClient) GetServicePorts(ctx context.Context, serviceName string, namespace string) ([]ServicePort, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Return mock ports for testing
	return []ServicePort{
		{Name: "http", Port: 80, TargetPort: 8080, Protocol: "TCP"},
		{Name: "https", Port: 443, TargetPort: 8443, Protocol: "TCP"},
	}, nil
}

func (m *mockK8sClient) PortForward(ctx context.Context, serviceName string, namespace string, localPort int, servicePort int32) error {
	// Mock implementation - just return the error if any
	return m.err
}

// mockNgrokClient implements a mock ngrok client for testing
type mockNgrokClient struct {
	startTunnelError error
	closeError       error
}

func (m *mockNgrokClient) StartTunnel(ctx context.Context, port int) (string, error) {
	if m.startTunnelError != nil {
		return "", m.startTunnelError
	}
	return fmt.Sprintf("https://mock%d.ngrok.io", port), nil
}

func (m *mockNgrokClient) Close() error {
	return m.closeError
}

func TestNewMockService(t *testing.T) {
	mockClient := &mockK8sClient{}
	mockNgrok := &mockNgrokClient{}
	svc := NewService(mockClient, mockNgrok)
	if svc == nil {
		t.Fatal("NewService should return a non-nil service")
	}
}

func TestGetServices(t *testing.T) {
	expectedServices := []string{
		"web-frontend",
		"api-gateway",
		"user-service",
		"database-service",
		"cache-service",
		"notification-service",
	}

	mockClient := &mockK8sClient{
		services: expectedServices,
	}
	mockNgrok := &mockNgrokClient{}
	svc := NewService(mockClient, mockNgrok)
	services, err := svc.GetServices(context.Background())

	if err != nil {
		t.Fatalf("GetServices should not return an error: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("GetServices should return at least one service")
	}

	if len(services) != len(expectedServices) {
		t.Fatalf("Expected %d services, got %d", len(expectedServices), len(services))
	}

	for i, expected := range expectedServices {
		if services[i] != expected {
			t.Errorf("Expected service %s at index %d, got %s", expected, i, services[i])
		}
	}
}

func TestGetServicesWithNilClient(t *testing.T) {
	mockNgrok := &mockNgrokClient{}
	svc := NewService(nil, mockNgrok)
	services, err := svc.GetServices(context.Background())

	if err == nil {
		t.Fatal("GetServices should return an error when client is nil")
	}

	if services != nil {
		t.Error("GetServices should return nil services when client is nil")
	}
}

func TestGetServicePorts(t *testing.T) {
	mockClient := &mockK8sClient{}
	mockNgrok := &mockNgrokClient{}
	svc := NewService(mockClient, mockNgrok)

	ports, err := svc.GetServicePorts(context.Background(), "test-service (ns: default)")
	if err != nil {
		t.Fatalf("GetServicePorts should not return an error: %v", err)
	}

	if len(ports) != 2 {
		t.Fatalf("Expected 2 ports, got %d", len(ports))
	}

	// Check first port
	if ports[0].Name != "http" || ports[0].Port != 80 || ports[0].TargetPort != 8080 {
		t.Errorf("First port should be http:80->8080, got %s:%d->%d", ports[0].Name, ports[0].Port, ports[0].TargetPort)
	}

	// Check second port
	if ports[1].Name != "https" || ports[1].Port != 443 || ports[1].TargetPort != 8443 {
		t.Errorf("Second port should be https:443->8443, got %s:%d->%d", ports[1].Name, ports[1].Port, ports[1].TargetPort)
	}
}

func TestStartPortForwarding(t *testing.T) {
	mockClient := &mockK8sClient{}
	mockNgrok := &mockNgrokClient{}
	svc := NewService(mockClient, mockNgrok)
	// Use the proper format with namespace
	port, err := svc.StartPortForwarding(context.Background(), "test-service (ns: default)", 80)

	if err != nil {
		t.Fatalf("StartPortForwarding should not return an error: %v", err)
	}

	if port < 8000 || port >= 9000 {
		t.Errorf("Port should be between 8000-8999, got %d", port)
	}
}

func TestCreateNgrokSession(t *testing.T) {
	mockClient := &mockK8sClient{}
	mockNgrok := &mockNgrokClient{}
	svc := NewService(mockClient, mockNgrok)
	url, err := svc.CreateNgrokSession(context.Background(), 8080)

	if err != nil {
		t.Fatalf("CreateNgrokSession should not return an error: %v", err)
	}

	if url == "" {
		t.Fatal("CreateNgrokSession should return a non-empty URL")
	}

	// Check if URL has expected format
	if len(url) < 10 || url[:8] != "https://" || url[len(url)-9:] != ".ngrok.io" {
		t.Errorf("URL should have format https://xxxxx.ngrok.io, got %s", url)
	}
}

func TestCleanup(t *testing.T) {
	mockClient := &mockK8sClient{}
	mockNgrok := &mockNgrokClient{}
	svc := NewService(mockClient, mockNgrok)

	// Start some services to cleanup
	_, err := svc.StartPortForwarding(context.Background(), "test-service (ns: default)", 80)
	if err != nil {
		t.Fatalf("StartPortForwarding should not return an error: %v", err)
	}
	_, err = svc.CreateNgrokSession(context.Background(), 8080)
	if err != nil {
		t.Fatalf("CreateNgrokSession should not return an error: %v", err)
	}
	// Test cleanup
	err = svc.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup should not return an error: %v", err)
	}
	// Verify cleanup cleared the state
	if svc.activeService != "" {
		t.Error("activeService should be empty after cleanup")
	}

	if svc.activePort != 0 {
		t.Error("activePort should be 0 after cleanup")
	}

	if svc.activeNgrokURL != "" {
		t.Error("activeNgrokURL should be empty after cleanup")
	}
}
