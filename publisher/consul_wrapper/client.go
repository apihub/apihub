package consul_wrapper

import "github.com/hashicorp/consul/api"

type Client struct {
	*api.Client
	runner *Runner
}

func NewClient(runner *Runner) (*Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: runner.Address,
		Scheme:  runner.Scheme,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: client,
	}, nil
}
