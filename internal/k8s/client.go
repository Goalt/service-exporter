package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type client struct {
	clientset *kubernetes.Clientset
}

func New() (*client, func()) {
	// Get kubeconfig path from environment variable
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		// Fall back to default kubeconfig location
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Warning: could not get home directory: %v\n", err)
			return nil, func() {}
		}
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Warning: could not build config from kubeconfig: %v\n", err)
		return nil, func() {}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Warning: could not create kubernetes client: %v\n", err)
		return nil, func() {}
	}

	return &client{clientset: clientset}, func() {}
}

func (c *client) ListServices() ([]string, error) {
	if c.clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	// List services in all namespaces
	services, err := c.clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var serviceNames []string
	for _, svc := range services.Items {
		// Include namespace in the service name for clarity
		serviceName := fmt.Sprintf("%s (ns: %s)", svc.Name, svc.Namespace)
		serviceNames = append(serviceNames, serviceName)
	}

	return serviceNames, nil
}

func (c *client) PortForward(serviceName string, namespace string, port int) error {
	// Port forwarding implementation would go here
	// For simplicity, we will just simulate success
	return nil
}
