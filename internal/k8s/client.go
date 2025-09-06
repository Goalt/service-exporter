package k8s

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

	"github.com/Goalt/service-exporter/internal/service"
)

type client struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

func New(kubeconfigPath string) (*client, error) {
	// Use provided kubeconfig path or fall back to default
	if kubeconfigPath == "" {
		// Fall back to default kubeconfig location
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("kubeconfig path not provided and could not determine home directory: %w", err)
		}

		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig from path %s: %w", kubeconfigPath, err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &client{clientset: clientset, config: config}, nil
}

func (c *client) ListServices(ctx context.Context) ([]string, error) {
	if c.clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	// List services in all namespaces
	services, err := c.clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
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

func (c *client) GetServicePorts(ctx context.Context, serviceName string, namespace string) ([]service.ServicePort, error) {
	if c.clientset == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	// Get the service to find available ports
	svc, err := c.clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s in namespace %s: %w", serviceName, namespace, err)
	}

	if len(svc.Spec.Ports) == 0 {
		return nil, fmt.Errorf("service %s has no ports defined", serviceName)
	}

	var servicePorts []service.ServicePort
	for _, port := range svc.Spec.Ports {
		targetPort := port.TargetPort.IntVal
		if targetPort == 0 {
			// If TargetPort is not specified, use the service port
			targetPort = port.Port
		}

		servicePorts = append(servicePorts, service.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: targetPort,
			Protocol:   string(port.Protocol),
		})
	}

	return servicePorts, nil
}

func (c *client) PortForward(ctx context.Context, serviceName string, namespace string, localPort int, servicePort int32) error {
	if c.clientset == nil || c.config == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	// Get the service to find target port
	svc, err := c.clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get service %s in namespace %s: %w", serviceName, namespace, err)
	}

	if len(svc.Spec.Ports) == 0 {
		return fmt.Errorf("service %s has no ports defined", serviceName)
	}

	// Find the specific port that matches the requested servicePort
	var selectedPort *corev1.ServicePort
	for _, port := range svc.Spec.Ports {
		if port.Port == servicePort {
			selectedPort = &port
			break
		}
	}

	if selectedPort == nil {
		return fmt.Errorf("port %d not found in service %s", servicePort, serviceName)
	}

	// Find pods that match the service selector
	pods, err := c.findPodsForService(ctx, svc)
	if err != nil {
		return fmt.Errorf("failed to find pods for service %s: %w", serviceName, err)
	}

	if len(pods) == 0 {
		return fmt.Errorf("no running pods found for service %s", serviceName)
	}

	// Use the first available pod
	pod := pods[0]

	// Determine the target port on the pod
	targetPort := selectedPort.TargetPort.IntVal
	if targetPort == 0 {
		// If TargetPort is not specified, use the service port
		targetPort = selectedPort.Port
	}

	// Create port forward request
	req := c.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(pod.Name).
		SubResource("portforward")

	// Create SPDY transport
	transport, upgrader, err := spdy.RoundTripperFor(c.config)
	if err != nil {
		return fmt.Errorf("failed to create SPDY transport: %w", err)
	}

	// Create dialer
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())

	// Set up port forwarding
	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})

	ports := []string{fmt.Sprintf("%d:%d", localPort, targetPort)}

	pf, err := portforward.New(dialer, ports, stopCh, readyCh, io.Discard, io.Discard)
	if err != nil {
		return fmt.Errorf("failed to create port forwarder: %w", err)
	}

	// Start port forwarding in a goroutine
	go func() {
		if err := pf.ForwardPorts(); err != nil {
			log.Printf("Port forwarding error: %v\n", err)
		}
	}()

	// Wait for port forwarding to be ready or timeout
	select {
	case <-readyCh:
		log.Printf("Port forwarding ready from localhost:%d to pod %s:%d\n", localPort, pod.Name, targetPort)
		return nil
	case <-time.After(30 * time.Second):
		close(stopCh)
		return fmt.Errorf("timeout waiting for port forwarding to be ready")
	}
}

func (c *client) findPodsForService(ctx context.Context, svc *corev1.Service) ([]corev1.Pod, error) {
	// Convert service selector to label selector string
	labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: svc.Spec.Selector})

	pods, err := c.clientset.CoreV1().Pods(svc.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}

	// Filter only running pods
	var runningPods []corev1.Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			runningPods = append(runningPods, pod)
		}
	}

	return runningPods, nil
}
