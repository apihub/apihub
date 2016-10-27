package client

import (
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client/connection"
)

type client struct {
	conn connection.Connection
}

func New(conn connection.Connection) apihub.Client {
	return &client{
		conn: conn,
	}
}

func (cli *client) Ping() error {
	return cli.conn.Ping()
}

func (cli *client) AddService(spec apihub.ServiceSpec) (apihub.Service, error) {
	service, err := cli.conn.AddService(spec)
	if err != nil {
		return nil, err
	}

	return newService(service.Handle, cli.conn), nil
}

func (cli *client) Services() ([]apihub.Service, error) {
	specs, err := cli.conn.Services()
	if err != nil {
		return nil, err
	}

	services := []apihub.Service{}
	for _, spec := range specs {
		services = append(services, newService(spec.Handle, cli.conn))
	}

	return services, nil
}

func (cli *client) RemoveService(handle string) error {
	return nil
}

func (cli *client) Lookup(handle string) (apihub.Service, error) {
	return nil, nil
}
