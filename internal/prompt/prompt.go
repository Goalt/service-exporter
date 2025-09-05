package prompt

import (
	"errors"
	"fmt"
	"strconv"

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
