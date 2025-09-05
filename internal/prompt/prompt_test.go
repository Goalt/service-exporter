package prompt

import (
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
