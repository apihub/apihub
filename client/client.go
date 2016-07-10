package client

import (
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client/connection"
)

type client struct {
	connection connection.Connection
}

func New(connection connection.Connection) apihub.Client {
	return &client{
		connection: connection,
	}
}

func (cli *client) Ping() error {
	return nil
}

func (cli *client) AddService(apihub.ServiceSpec) (apihub.Service, error) {
	return nil, nil
}

func (cli *client) RemoveService(handle string) error {
	return nil
}

func (cli *client) Services() ([]apihub.Service, error) {
	return nil, nil
}

func (cli *client) Lookup(handle string) (apihub.Service, error) {
	return nil, nil
}

func (cli *client) do(handle string) error {
	return nil
}
