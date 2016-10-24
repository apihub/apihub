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

func (s *service) Backends() []apihub.Backend {
	panic("not implemented")
}

func (s *service) Start() error {
	panic("not implemented")
}

func (s *service) Stop(kill bool) error {
	panic("not implemented")
}

func (s *service) Info() (apihub.ServiceSpec, error) {
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

func (s *service) Timeout() time.Duration {
	panic("not implemented")
}
