package ngrok

import (
	"context"
	"testing"
	"time"
)

func TestClient_StartTunnel(t *testing.T) {
	// This test requires an ngrok auth token to work
	// We'll create a basic unit test that checks the interface

	// Test that we can create a client with empty auth token (should fail)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := NewClient(ctx, "")
	if err == nil {
		t.Error("Expected error when creating client with empty auth token")
	}
}

func TestClient_StartTunnel_InvalidToken(t *testing.T) {
	// Test with invalid auth token
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient(ctx, "invalid_token")
	if err != nil {
		t.Errorf("NewClient should not fail with invalid token, got: %v", err)
		return
	}

	// The error should occur when trying to start a tunnel
	_, err = client.StartTunnel(ctx, 8080)
	if err == nil {
		t.Error("Expected error when starting tunnel with invalid auth token")
	}
}

func TestClient_Close(t *testing.T) {
	// Test that we can call Close on a nil session client
	client := &Client{}
	err := client.Close()
	if err != nil {
		t.Errorf("Close() should not return error for nil session, got: %v", err)
	}
}
