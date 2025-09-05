package ngrok

type client struct{}

func New() *client {
	return &client{}
}

func (c *client) StartTunnel(port int) (string, error) {
	// Placeholder implementation
	return "https://example.ngrok.io", nil
}
