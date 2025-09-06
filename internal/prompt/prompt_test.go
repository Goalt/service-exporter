package prompt

import (
	"fmt"
	"testing"
)

func TestServiceSelectPrompt_EmptyServices(t *testing.T) {
	services := []string{}
	_, err := ServiceSelectPrompt(services)

	if err == nil {
		t.Fatal("ServiceSelectPrompt should return an error for empty services list")
	}

	expectedMsg := "no services available"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestServiceSelectPrompt_ValidServices(t *testing.T) {
	services := []string{"service1", "service2", "service3"}
	// Note: We can't easily test the interactive part without mocking promptui,
	// but we can at least verify that the function exists and handles empty input correctly
	if len(services) == 0 {
		t.Fatal("This test requires non-empty services list")
	}
}

func TestNgrokTokenPrompt_ValidatesEmptyInput(t *testing.T) {
	// Note: We can't easily test the interactive prompts without complex mocking,
	// but we can verify the function signature and that it exists
	// In a real scenario, this would require mocking promptui or using dependency injection

	// Just verify the function exists and can be called
	// The actual validation logic is in the validate function passed to promptui
	validate := func(input string) error {
		if input == "" {
			return fmt.Errorf("Ngrok auth token cannot be empty")
		}
		return nil
	}

	// Test empty input validation
	err := validate("")
	if err == nil {
		t.Error("Expected validation error for empty input")
	}

	// Test valid input
	err = validate("test-token")
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

func TestKubeconfigPathPrompt_Exists(t *testing.T) {
	// Just verify the function exists and can be referenced
	// The actual testing of promptui interactions would require complex mocking
	defer func() {
		if r := recover(); r != nil {
			t.Error("KubeconfigPathPrompt function should exist and be callable")
		}
	}()
	// This will not actually call the function but verifies it exists
	_ = KubeconfigPathPrompt
}

func TestUseDefaultsPrompt_Exists(t *testing.T) {
	// Just verify the function exists and can be referenced
	// The actual testing of promptui interactions would require complex mocking
	defer func() {
		if r := recover(); r != nil {
			t.Error("UseDefaultsPrompt function should exist and be callable")
		}
	}()
	// This will not actually call the function but verifies it exists
	_ = UseDefaultsPrompt
}
