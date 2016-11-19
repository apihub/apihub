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

	return newService(service.Host, cli.conn), nil
}

func (cli *client) Services() ([]apihub.Service, error) {
	specs, err := cli.conn.Services()
	if err != nil {
		return nil, err
	}

	services := []apihub.Service{}
	for _, spec := range specs {
		services = append(services, newService(spec.Host, cli.conn))
	}

	return services, nil
}

func (cli *client) RemoveService(host string) error {
	return cli.conn.RemoveService(host)
}

func (cli *client) FindService(host string) (apihub.Service, error) {
	service, err := cli.conn.FindService(host)
	if err != nil {
		return nil, err
	}

	return newService(service.Host, cli.conn), nil
}

func (cli *client) UpdateService(host string, spec apihub.ServiceSpec) (apihub.Service, error) {
	service, err := cli.conn.UpdateService(host, spec)
	if err != nil {
		return nil, err
	}

	return newService(service.Host, cli.conn), nil
}
