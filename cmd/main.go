package main

import (
	"fmt"
	"log"

	"github.com/Goalt/service-exporter/internal/prompt"
	"github.com/Goalt/service-exporter/internal/service"
)

func main() {
	fmt.Println("🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok")
	fmt.Println("================================================================")

	// Initialize the service
	svc := service.NewMockService()

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
}
