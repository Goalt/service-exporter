package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goalt/service-exporter/internal/k8s"
	"github.com/Goalt/service-exporter/internal/ngrok"
	"github.com/Goalt/service-exporter/internal/prompt"
	"github.com/Goalt/service-exporter/internal/service"
)

// Config holds all environment configuration
type Config struct {
	NgrokAuthToken string
	KubeconfigPath string
}

// loadConfig reads configuration from environment variables or prompts user for input
func loadConfig() (*Config, error) {
	log.Println("\n⚙️  Configuration Setup")
	log.Println("=====================")

	// Ask user if they want to use defaults or provide manual input
	useDefaults, err := prompt.UseDefaultsPrompt()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration preference: %v", err)
	}

	config := &Config{}

	if useDefaults {
		log.Println("\n📋 Using environment variables for configuration...")
		config.NgrokAuthToken = os.Getenv("NGROK_AUTH_TOKEN")
		config.KubeconfigPath = os.Getenv("KUBECONFIG")

		// Validate required environment variables when using defaults
		if config.NgrokAuthToken == "" {
			return nil, fmt.Errorf("❌ NGROK_AUTH_TOKEN environment variable is required when using default configuration")
		}
	} else {
		log.Println("\n📝 Manual configuration mode...")

		// Prompt for ngrok auth token
		config.NgrokAuthToken, err = prompt.NgrokTokenPrompt()
		if err != nil {
			return nil, fmt.Errorf("failed to get ngrok auth token: %v", err)
		}

		// Prompt for kubeconfig path
		config.KubeconfigPath, err = prompt.KubeconfigPathPrompt()
		if err != nil {
			return nil, fmt.Errorf("failed to get kubeconfig path: %v", err)
		}
	}

	return config, nil
}

func main() {
	log.Println("🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok")
	log.Println("================================================================")

	// Load configuration from prompts or environment variables
	config, err := loadConfig()
	if err != nil {
		log.Printf("❌ Configuration failed: %v", err)
		return
	}

	// create Kubernetes client with kubeconfig path
	k8sClient, err := k8s.New(config.KubeconfigPath)
	if err != nil {
		log.Printf("❌ Failed to create Kubernetes client: %v", err)
		return
	}

	log.Println("🔑 Found ngrok auth token, creating ngrok client")
	ngrokClient, err := ngrok.NewClient(config.NgrokAuthToken)
	if err != nil {
		log.Printf("❌ Failed to create ngrok client: %v", err)
		return
	}

	svc := service.NewService(k8sClient, ngrokClient)

	// Setup graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// create a context that is canceled on interrupt signal
	ctx, cancel := context.WithCancel(context.Background())

	// Setup cleanup function
	cleanup := func() {
		cancel()
		if err := svc.Cleanup(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
		log.Println("\n👋 Goodbye!")
		os.Exit(0)
	}

	// Handle signals in a goroutine
	go func() {
		<-sigChan
		cleanup()
	}()

	// Step 1: Get list of Kubernetes services
	log.Println("\n📋 Fetching available Kubernetes services...")
	k8sServices, err := svc.GetServices(ctx)
	if err != nil {
		log.Printf("Failed to get services: %v", err)
		return
	}

	// Step 2: User selects a service
	selectedK8SService, err := prompt.ServiceSelectPrompt(k8sServices)
	if err != nil {
		log.Printf("Service selection failed: %v", err)
		return
	}

	log.Printf("\n✅ Selected service: %s\n", selectedK8SService)

	// Step 3: Start port forwarding
	port, err := svc.StartPortForwarding(ctx, selectedK8SService)
	if err != nil {
		log.Printf("Failed to start port forwarding: %v", err)
		return
	}

	// Step 4: Create ngrok session
	ngrokURL, err := svc.CreateNgrokSession(ctx, port)
	if err != nil {
		log.Printf("Failed to create ngrok session: %v", err)
		return
	}

	// Display final result
	log.Println("\n🎉 Setup complete!")
	log.Println("==================")
	log.Printf("Service: %s\n", selectedK8SService)
	log.Printf("Local Port: %d\n", port)
	log.Printf("Public URL: %s\n", ngrokURL)
	log.Println("\nYou can now access your service via the public URL above!")
	log.Println("\n📌 Press Ctrl+C to gracefully shutdown and cleanup resources...")

	// Keep the program running until interrupted
	select {}
}
