package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Goalt/service-exporter/internal/k8s"
	"github.com/Goalt/service-exporter/internal/ngrok"
	"github.com/Goalt/service-exporter/internal/prompt"
	"github.com/Goalt/service-exporter/internal/service"
)

type App struct {
	config Config
	svc    service.Service
}

func New() *App {
	return &App{}
}

func (a *App) LoadConfig() error {
	log.Println("🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok")
	log.Println("================================================================")

	// Load configuration from prompts or environment variables
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	a.config = config
	return nil
}

func (a *App) Run(ctx context.Context) error {
	// create Kubernetes client with kubeconfig path
	k8sClient, err := k8s.New(a.config.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	log.Println("🔑 Found ngrok auth token, creating ngrok client")
	ngrokClient, err := ngrok.NewClient(a.config.NgrokAuthToken)
	if err != nil {
		return fmt.Errorf("failed to create ngrok client: %v", err)
	}

	a.svc = service.NewService(k8sClient, ngrokClient)

	// Step 1: Get list of Kubernetes services
	log.Println("\n📋 Fetching available Kubernetes services...")
	k8sServices, err := a.svc.GetServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to get services: %v", err)
	}

	// Step 2: User selects a service
	selectedK8SService, err := prompt.ServiceSelectPrompt(k8sServices)
	if err != nil {
		return fmt.Errorf("service selection failed: %v", err)
	}

	log.Printf("\n✅ Selected service: %s\n", selectedK8SService)

	// Step 3: Get available ports for the selected service
	log.Println("\n📋 Fetching available ports for the selected service...")
	servicePorts, err := a.svc.GetServicePorts(ctx, selectedK8SService)
	if err != nil {
		return fmt.Errorf("failed to get service ports: %v", err)
	}

	// Step 4: User selects a port to forward
	selectedPort, err := prompt.PortSelectPrompt(servicePorts)
	if err != nil {
		return fmt.Errorf("port selection failed: %v", err)
	}

	log.Printf("\n✅ Selected port: %d (%s)\n", selectedPort.Port, selectedPort.Name)

	// Step 5: Start port forwarding
	port, err := a.svc.StartPortForwarding(ctx, selectedK8SService, selectedPort.Port)
	if err != nil {
		return fmt.Errorf("failed to start port forwarding: %v", err)
	}

	// Step 6: Create ngrok session
	ngrokURL, err := a.svc.CreateNgrokSession(ctx, port)
	if err != nil {
		return fmt.Errorf("failed to create ngrok session: %v", err)
	}

	// Display final result
	log.Println("\n🎉 Setup complete!")
	log.Println("==================")
	log.Printf("Service: %s\n", selectedK8SService)
	portName := selectedPort.Name
	if portName == "" {
		portName = "unnamed"
	}
	log.Printf("Selected Port: %d (%s)\n", selectedPort.Port, portName)
	log.Printf("Local Port: %d\n", port)
	log.Printf("Public URL: %s\n", ngrokURL)
	log.Println("\nYou can now access your service via the public URL above!")
	log.Println("\n📌 Press Ctrl+C to gracefully shutdown and cleanup resources...")

	<-ctx.Done()

	return nil
}

func (a *App) Cleanup() error {
	if err := a.svc.Cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup resources: %v", err)
	}

	log.Println("\n👋 Goodbye!")

	return nil
}
