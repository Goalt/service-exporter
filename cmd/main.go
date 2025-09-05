package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goalt/service-exporter/internal/prompt"
	"github.com/Goalt/service-exporter/internal/service"
)

func main() {
	fmt.Println("🚀 Service Exporter - Kubernetes Service Port Forwarding with ngrok")
	fmt.Println("================================================================")

	// Initialize the service
	svc := service.NewMockService()

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
