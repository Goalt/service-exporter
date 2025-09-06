package ngrok

import (
	"context"
	"fmt"
	"net/url"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

// Client represents an ngrok client for creating tunnels
type Client struct {
	authToken string
	forwarder ngrok.Forwarder
}

// NewClient creates a new ngrok client
func NewClient(authToken string) (*Client, error) {
	if authToken == "" {
		return nil, fmt.Errorf("auth token is required")
	}

	return &Client{
		authToken: authToken,
	}, nil
}

// StartTunnel creates a new HTTP tunnel for the specified port
func (c *Client) StartTunnel(ctx context.Context, port int) (string, error) {
	// Create backend URL
	backendURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		return "", fmt.Errorf("failed to parse backend URL: %w", err)
	}

	// Use the simplified ListenAndForward function which handles everything
	forwarder, err := ngrok.ListenAndForward(ctx, backendURL, config.HTTPEndpoint(), ngrok.WithAuthtoken(c.authToken))
	if err != nil {
		return "", fmt.Errorf("failed to create tunnel: %w", err)
	}

	// Store forwarder for cleanup
	c.forwarder = forwarder

	return forwarder.URL(), nil
}

// Close closes the ngrok forwarder
func (c *Client) Close() error {
	if c.forwarder != nil {
		return c.forwarder.Close()
	}
	return nil
}
