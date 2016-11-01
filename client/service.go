package client

import (
	"time"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client/connection"
)

type service struct {
	handle string
	conn   connection.Connection
}

func newService(handle string, conn connection.Connection) apihub.Service {
	return &service{
		handle: handle,
		conn:   conn,
	}
}

func (s *service) Handle() string {
	return s.handle
}

func (s *service) Info() (apihub.ServiceSpec, error) {
	return s.conn.FindService(s.Handle())
}

func (s *service) Backends() ([]apihub.BackendInfo, error) {
	info, err := s.conn.FindService(s.Handle())
	if err != nil {
		return nil, err
	}

	return info.Backends, nil
}

func (s *service) Start() error {
	panic("not implemented")
}

func (s *service) Stop(kill bool) error {
	panic("not implemented")
}

func (s *service) AddBackend(be apihub.BackendInfo) error {
	panic("not implemented")
}

func (s *service) RemoveBackend(be apihub.BackendInfo) error {
	panic("not implemented")
}

func (s *service) SetTimeout(time.Duration) {
	panic("not implemented")
}
