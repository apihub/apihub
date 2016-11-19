package client

import (
	"time"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client/connection"
)

type service struct {
	host string
	conn   connection.Connection
}

func newService(host string, conn connection.Connection) apihub.Service {
	return &service{
		host: host,
		conn:   conn,
	}
}

func (s *service) Host() string {
	return s.host
}

func (s *service) Info() (apihub.ServiceSpec, error) {
	return s.conn.FindService(s.Host())
}

func (s *service) Start() error {
	spec := apihub.ServiceSpec{
		Disabled: false,
	}
	_, err := s.conn.UpdateService(s.Host(), spec)
	return err
}

func (s *service) Stop() error {
	spec := apihub.ServiceSpec{
		Disabled: true,
	}
	_, err := s.conn.UpdateService(s.Host(), spec)
	return err
}

func (s *service) SetTimeout(duration time.Duration) error {
	spec := apihub.ServiceSpec{
		Timeout: duration,
	}
	_, err := s.conn.UpdateService(s.Host(), spec)
	return err
}

func (s *service) Backends() ([]apihub.BackendInfo, error) {
	info, err := s.conn.FindService(s.Host())
	if err != nil {
		return nil, err
	}

	return info.Backends, nil
}
