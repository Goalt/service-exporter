package ngrok

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

// Client represents an ngrok client for creating tunnels
type Client struct {
	session ngrok.Session
}

// NewClient creates a new ngrok client
func NewClient(ctx context.Context, authToken string) (*Client, error) {
	sessionCh := make(chan ngrok.Session)
	errCh := make(chan error)
	go func() {
		session, err := ngrok.Connect(ctx, ngrok.WithAuthtoken(authToken))
		if err != nil {
			errCh <- err
			return
		}
		sessionCh <- session
	}()

	var session ngrok.Session
	select {
	case session = <-sessionCh:
		// Successfully connected
	case err := <-errCh:
		return nil, fmt.Errorf("failed to connect to ngrok: %w", err)
	}

	return &Client{
		session: session,
	}, nil
}

// StartTunnel creates a new HTTP tunnel for the specified port
func (c *Client) StartTunnel(ctx context.Context, port int) (string, error) {
	tunnel, err := c.session.Listen(ctx, config.HTTPEndpoint())
	if err != nil {
		return "", fmt.Errorf("failed to create tunnel: %w", err)
	}

	// Create a reverse proxy to forward requests to localhost:port
	targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		tunnel.Close()
		return "", fmt.Errorf("failed to parse target URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Start serving the tunnel
	go func() {
		defer tunnel.Close()
		err := http.Serve(tunnel, proxy)
		if err != nil && err != http.ErrServerClosed {
			// Log error if needed, but don't block
			fmt.Printf("Error serving tunnel: %v\n", err)
		}
	}()

	return tunnel.URL(), nil
}

// Close closes the ngrok session
func (c *Client) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}
