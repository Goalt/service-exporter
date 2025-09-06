package app

import (
	"fmt"
	"log"
	"os"

	"github.com/Goalt/service-exporter/internal/prompt"
)

// Config holds all environment configuration
type Config struct {
	NgrokAuthToken string
	KubeconfigPath string
}

// loadConfig reads configuration from environment variables or prompts user for input
func loadConfig() (Config, error) {
	log.Println("\n⚙️  Configuration Setup")
	log.Println("=====================")

	// Ask user if they want to use defaults or provide manual input
	useDefaults, err := prompt.UseDefaultsPrompt()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get configuration preference: %v", err)
	}

	config := Config{}

	if useDefaults {
		log.Println("\n📋 Using environment variables for configuration...")
		config.NgrokAuthToken = os.Getenv("NGROK_AUTH_TOKEN")
		config.KubeconfigPath = os.Getenv("KUBECONFIG")

		// Validate required environment variables when using defaults
		if config.NgrokAuthToken == "" {
			return Config{}, fmt.Errorf("❌ NGROK_AUTH_TOKEN environment variable is required when using default configuration")
		}
	} else {
		log.Println("\n📝 Manual configuration mode...")

		// Prompt for ngrok auth token
		config.NgrokAuthToken, err = prompt.NgrokTokenPrompt()
		if err != nil {
			return Config{}, fmt.Errorf("failed to get ngrok auth token: %v", err)
		}

		// Prompt for kubeconfig path
		config.KubeconfigPath, err = prompt.KubeconfigPathPrompt()
		if err != nil {
			return Config{}, fmt.Errorf("failed to get kubeconfig path: %v", err)
		}
	}

	return config, nil
}
