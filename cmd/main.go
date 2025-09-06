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
	fmt.Println("\n⚙️  Configuration Setup")
	fmt.Println("=====================")
	
	// Ask user if they want to use defaults or provide manual input
	useDefaults, err := prompt.UseDefaultsPrompt()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration preference: %v", err)
	}

	config := &Config{}

	if useDefaults {
		fmt.Println("\n📋 Using environment variables for configuration...")
		config.NgrokAuthToken = os.Getenv("NGROK_AUTH_TOKEN")
		config.KubeconfigPath = os.Getenv("KUBECONFIG")
		
		// Validate required environment variables when using defaults
		if config.NgrokAuthToken == "" {
			return nil, fmt.Errorf("❌ NGROK_AUTH_TOKEN environment variable is required when using default configuration")
		}
	} else {
		fmt.Println("\n📝 Manual configuration mode...")
		
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
	fmt.Println("🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok")
	fmt.Println("================================================================")

	// Load configuration from prompts or environment variables
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("❌ Configuration failed: %v", err)
	}

	// create Kubernetes client with kubeconfig path
	k8sClient, k8sCleanup := k8s.New(config.KubeconfigPath)
	defer k8sCleanup()

	fmt.Println("🔑 Found ngrok auth token, creating ngrok client")
	ngrokClient, err := ngrok.NewClient(context.Background(), config.NgrokAuthToken)
	if err != nil {
		log.Fatalf("❌ Failed to create ngrok client: %v", err)
	}

	svc := service.NewService(k8sClient, ngrokClient)

	// Setup graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Setup cleanup function
	cleanup := func() {
		if err := svc.Cleanup(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
		fmt.Println("\n👋 Goodbye!")
		os.Exit(0)
	}

	// Handle signals in a goroutine
	go func() {
		<-sigChan
		cleanup()
	}()

	// Step 1: Get list of Kubernetes services
	fmt.Println("\n📋 Fetching available Kubernetes services...")
	services, err := svc.GetServices()
	if err != nil {
		log.Fatalf("Failed to get services: %v", err)
	}

	// Step 2: User selects a service
	selectedService, err := prompt.ServiceSelectPrompt(services)
	if err != nil {
		log.Fatalf("Service selection failed: %v", err)
	}

	fmt.Printf("\n✅ Selected service: %s\n", selectedService)

	// Step 3: Start port forwarding
	port, err := svc.StartPortForwarding(selectedService)
	if err != nil {
		log.Fatalf("Failed to start port forwarding: %v", err)
	}

	// Step 4: Create ngrok session
	ngrokURL, err := svc.CreateNgrokSession(port)
	if err != nil {
		log.Fatalf("Failed to create ngrok session: %v", err)
	}

	// Display final result
	fmt.Println("\n🎉 Setup complete!")
	fmt.Println("==================")
	fmt.Printf("Service: %s\n", selectedService)
	fmt.Printf("Local Port: %d\n", port)
	fmt.Printf("Public URL: %s\n", ngrokURL)
	fmt.Println("\nYou can now access your service via the public URL above!")
	fmt.Println("\n📌 Press Ctrl+C to gracefully shutdown and cleanup resources...")

	// Keep the program running until interrupted
	select {}
}
