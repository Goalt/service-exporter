package k8s

type client struct{}

func New() (*client, func()) {
	// returns a new Kubernetes client and a cleanup function
	return &client{}, func() {}
}

func (c *client) ListServices() ([]string, error) {
	// Implementation to list services from Kubernetes cluster
	return nil, nil
}
