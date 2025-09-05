package service

import (
	"testing"
)

func TestNewMockService(t *testing.T) {
	svc := NewService(nil)
	if svc == nil {
		t.Fatal("NewMockService should return a non-nil service")
	}
}

func TestGetServices(t *testing.T) {
	svc := NewService(nil)
	services, err := svc.GetServices()

	if err != nil {
		t.Fatalf("GetServices should not return an error: %v", err)
	}

	if len(services) == 0 {
		t.Fatal("GetServices should return at least one service")
	}

	expectedServices := []string{
		"web-frontend",
		"api-gateway",
		"user-service",
		"database-service",
		"cache-service",
		"notification-service",
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

func TestStartPortForwarding(t *testing.T) {
	svc := NewService(nil)
	port, err := svc.StartPortForwarding("test-service")

	if err != nil {
		t.Fatalf("StartPortForwarding should not return an error: %v", err)
	}

	if port < 8000 || port >= 9000 {
		t.Errorf("Port should be between 8000-8999, got %d", port)
	}
}

func TestCreateNgrokSession(t *testing.T) {
	svc := NewService(nil)
	url, err := svc.CreateNgrokSession(8080)

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
	svc := NewService(nil).(*service)

	// Start some services to cleanup
	_, err := svc.StartPortForwarding("test-service")
	if err != nil {
		t.Fatalf("StartPortForwarding should not return an error: %v", err)
	}

	_, err = svc.CreateNgrokSession(8080)
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
