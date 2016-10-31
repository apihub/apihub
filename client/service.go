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

func (s *service) Backends() ([]apihub.Backend, error) {
	info, err := s.conn.FindService(s.Handle())
	if err != nil {
		return nil, err
	}

	var backends []apihub.Backend
	for _, backend := range info.Backends {
		backends = append(backends, newBackend(backend))
	}

	return backends, nil
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

func (s *service) Lookup(address string) (apihub.Backend, error) {
	panic("not implemented")
}

func (s *service) SetTimeout(time.Duration) {
	panic("not implemented")
}
