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
)

type client struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

func New(kubeconfigPath string) (*client, func(), error) {
	// Use provided kubeconfig path or fall back to default
	if kubeconfigPath == "" {
		// Fall back to default kubeconfig location
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Warning: could not get home directory: %v\n", err)
			return nil, func() {}, fmt.Errorf("kubeconfig path not provided and could not determine home directory")
		}
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Printf("Warning: could not build config from kubeconfig: %v\n", err)
		return nil, func() {}, fmt.Errorf("failed to build kubeconfig from path %s: %w", kubeconfigPath, err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: could not create kubernetes client: %v\n", err)
		return nil, func() {}, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &client{clientset: clientset, config: config}, func() {}, nil
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

func (c *client) PortForward(serviceName string, namespace string, localPort int) error {
	if c.clientset == nil || c.config == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	// Get the service to find target port
	svc, err := c.clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get service %s in namespace %s: %w", serviceName, namespace, err)
	}

	if len(svc.Spec.Ports) == 0 {
		return fmt.Errorf("service %s has no ports defined", serviceName)
	}

	// Use the first port of the service
	servicePort := svc.Spec.Ports[0]

	// Find pods that match the service selector
	pods, err := c.findPodsForService(svc)
	if err != nil {
		return fmt.Errorf("failed to find pods for service %s: %w", serviceName, err)
	}

	if len(pods) == 0 {
		return fmt.Errorf("no running pods found for service %s", serviceName)
	}

	// Use the first available pod
	pod := pods[0]

	// Determine the target port on the pod
	targetPort := servicePort.TargetPort.IntVal
	if targetPort == 0 {
		// If TargetPort is not specified, use the service port
		targetPort = servicePort.Port
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

func (c *client) findPodsForService(svc *corev1.Service) ([]corev1.Pod, error) {
	// Convert service selector to label selector string
	labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: svc.Spec.Selector})

	pods, err := c.clientset.CoreV1().Pods(svc.Namespace).List(context.TODO(), metav1.ListOptions{
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
