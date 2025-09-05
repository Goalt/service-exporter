package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
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
	activeService string
	activePort    int
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

// KubernetesService implements the Service interface with real Kubernetes functionality
type KubernetesService struct {
	clientset      *kubernetes.Clientset
	config         *rest.Config
	namespace      string
	activeService  string
	activePort     int
	activeNgrokURL string
	stopChannel    chan struct{}
	readyChannel   chan struct{}
}

// NewKubernetesService creates a new Kubernetes service instance
func NewKubernetesService(namespace string) (Service, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	return &KubernetesService{
		clientset: clientset,
		config:    config,
		namespace: namespace,
	}, nil
}

// getKubernetesConfig returns a Kubernetes client configuration
func getKubernetesConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if kubeconfigPath := os.Getenv("KUBECONFIG"); kubeconfigPath != "" {
		kubeconfig = kubeconfigPath
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	return config, nil
}

// GetServices returns a list of available Kubernetes services
func (k *KubernetesService) GetServices() ([]string, error) {
	ctx := context.Background()
	services, err := k.clientset.CoreV1().Services(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var serviceNames []string
	for _, service := range services.Items {
		serviceNames = append(serviceNames, service.Name)
	}

	return serviceNames, nil
}

// StartPortForwarding starts port forwarding for the specified service
func (k *KubernetesService) StartPortForwarding(serviceName string) (int, error) {
	ctx := context.Background()

	// Get the service to find its selector
	service, err := k.clientset.CoreV1().Services(k.namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	// Find pods that match the service selector
	selector := metav1.LabelSelector{MatchLabels: service.Spec.Selector}
	labelSelector := metav1.FormatLabelSelector(&selector)
	
	pods, err := k.clientset.CoreV1().Pods(k.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list pods for service %s: %w", serviceName, err)
	}

	if len(pods.Items) == 0 {
		return 0, fmt.Errorf("no pods found for service %s", serviceName)
	}

	// Use the first running pod
	var targetPod *corev1.Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			targetPod = &pod
			break
		}
	}

	if targetPod == nil {
		return 0, fmt.Errorf("no running pods found for service %s", serviceName)
	}

	// Get the first port from the service
	if len(service.Spec.Ports) == 0 {
		return 0, fmt.Errorf("service %s has no ports defined", serviceName)
	}

	targetPort := service.Spec.Ports[0].TargetPort.IntVal
	if targetPort == 0 {
		targetPort = service.Spec.Ports[0].Port
	}

	// Generate a random local port
	localPort := 8000 + rand.Intn(1000)

	// Setup port forwarding
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", k.namespace, targetPod.Name)
	hostIP := k.config.Host

	transport, upgrader, err := spdy.RoundTripperFor(k.config)
	if err != nil {
		return 0, fmt.Errorf("failed to create round tripper: %w", err)
	}

	serverURL, err := url.Parse(hostIP)
	if err != nil {
		return 0, fmt.Errorf("failed to parse host URL: %w", err)
	}
	serverURL.Path = path

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, serverURL)

	k.stopChannel = make(chan struct{}, 1)
	k.readyChannel = make(chan struct{}, 1)

	ports := []string{fmt.Sprintf("%d:%d", localPort, targetPort)}

	pf, err := portforward.New(dialer, ports, k.stopChannel, k.readyChannel, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create port forwarder: %w", err)
	}

	// Start port forwarding in a goroutine
	go func() {
		if err := pf.ForwardPorts(); err != nil {
			fmt.Printf("Error in port forwarding: %v\n", err)
		}
	}()

	// Wait for port forwarding to be ready
	select {
	case <-k.readyChannel:
		fmt.Printf("🔄 Port forwarding started for service '%s' on port %d -> %d\n", serviceName, localPort, targetPort)
	case <-time.After(30 * time.Second):
		close(k.stopChannel)
		return 0, fmt.Errorf("timeout waiting for port forwarding to be ready")
	}

	// Store the active session info
	k.activeService = serviceName
	k.activePort = localPort

	return localPort, nil
}

// CreateNgrokSession simulates creating an ngrok session
func (k *KubernetesService) CreateNgrokSession(port int) (string, error) {
	// Simulate some delay
	time.Sleep(2 * time.Second)
	
	// Generate a mock ngrok URL
	randomId := fmt.Sprintf("%x", rand.Uint32())
	ngrokURL := fmt.Sprintf("https://%s.ngrok.io", randomId)
	
	// Store the active ngrok URL
	k.activeNgrokURL = ngrokURL
	
	fmt.Printf("🌐 Creating ngrok tunnel for port %d...\n", port)
	
	return ngrokURL, nil
}

// Cleanup performs graceful shutdown of all active sessions
func (k *KubernetesService) Cleanup() error {
	fmt.Println("\n🔄 Performing graceful shutdown...")
	
	if k.activeNgrokURL != "" {
		fmt.Printf("🔌 Closing ngrok tunnel: %s\n", k.activeNgrokURL)
		time.Sleep(500 * time.Millisecond) // Simulate cleanup delay
		k.activeNgrokURL = ""
	}
	
	if k.stopChannel != nil {
		fmt.Printf("🔌 Stopping port forwarding for service '%s' on port %d\n", k.activeService, k.activePort)
		close(k.stopChannel)
		k.stopChannel = nil
		k.readyChannel = nil
		time.Sleep(500 * time.Millisecond) // Allow cleanup time
	}
	
	k.activeService = ""
	k.activePort = 0
	
	fmt.Println("✅ Graceful shutdown completed")
	return nil
}