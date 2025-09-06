package prompt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

// NumberPrompt prompts for a number input with validation
func NumberPrompt() (string, error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	return result, nil
}

// ServiceSelectPrompt prompts user to select a Kubernetes service
func ServiceSelectPrompt(services []string) (string, error) {
	if len(services) == 0 {
		return "", errors.New("no services available")
	}

	prompt := promptui.Select{
		Label: "Select a Kubernetes service",
		Items: services,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("service selection failed: %v", err)
	}

	return result, nil
}

// UseDefaultsPrompt asks user if they want to use default configuration or provide manual input
func UseDefaultsPrompt() (bool, error) {
	prompt := promptui.Select{
		Label: "Configuration mode",
		Items: []string{"Use default values from environment variables", "Provide parameters manually"},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("configuration mode selection failed: %v", err)
	}

	// Return true if user selected "Use default values" (index 0)
	return index == 0, nil
}

// NgrokTokenPrompt prompts user for ngrok auth token
func NgrokTokenPrompt() (string, error) {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return errors.New("Ngrok auth token cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Ngrok Auth Token",
		Validate: validate,
		Mask:     '*', // Hide the token input for security
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("ngrok token input failed: %v", err)
	}

	return strings.TrimSpace(result), nil
}

// KubeconfigPathPrompt prompts user for kubeconfig file path
func KubeconfigPathPrompt() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Kubeconfig file path (press Enter for default)",
		Default: "", // Empty default means it will use the system default
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("kubeconfig path input failed: %v", err)
	}

	return strings.TrimSpace(result), nil
}
