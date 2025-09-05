package k8s

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewClient_WithValidKubeconfig(t *testing.T) {
	// Create a temporary kubeconfig file for testing
	tmpDir, err := os.MkdirTemp("", "kubeconfig-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	kubeconfigPath := filepath.Join(tmpDir, "config")
	kubeconfigContent := `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://fake-k8s-server.example.com
  name: fake-cluster
contexts:
- context:
    cluster: fake-cluster
    user: fake-user
  name: fake-context
current-context: fake-context
users:
- name: fake-user
  user:
    token: fake-token
`

	err = os.WriteFile(kubeconfigPath, []byte(kubeconfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write kubeconfig: %v", err)
	}

	// Set environment variable
	os.Setenv("KUBECONFIG", kubeconfigPath)
	defer os.Unsetenv("KUBECONFIG")

	// Test client creation
	client, cleanup := New()
	defer cleanup()

	if client == nil {
		t.Error("Expected client to be created, got nil")
	}

	if client.clientset == nil {
		t.Error("Expected clientset to be initialized, got nil")
	}
}

func TestNewClient_WithoutKubeconfig(t *testing.T) {
	// Unset KUBECONFIG environment variable
	os.Unsetenv("KUBECONFIG")

	// Set a non-existent home directory to avoid finding real kubeconfig
	originalHome, homeExists := os.LookupEnv("HOME")
	if homeExists {
		defer os.Setenv("HOME", originalHome)
	}
	os.Setenv("HOME", "/non-existent-path")

	// Test client creation
	client, cleanup := New()
	defer cleanup()

	// Should return nil client when no valid kubeconfig is found
	if client != nil {
		t.Error("Expected client to be nil when no kubeconfig is available, got non-nil")
	}
}

func TestListServices_WithNilClient(t *testing.T) {
	client := &client{clientset: nil}

	services, err := client.ListServices()

	if err == nil {
		t.Error("Expected error when clientset is nil, got nil")
	}

	if services != nil {
		t.Error("Expected services to be nil when clientset is nil, got non-nil")
	}

	expectedError := "kubernetes client not initialized"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}
