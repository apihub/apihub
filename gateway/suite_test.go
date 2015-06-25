package gateway

import (
	"testing"

	"github.com/backstage/maestro/account"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	Settings *Settings
}

var service *account.Service

func (s *S) SetUpTest(c *C) {
	s.Settings = &Settings{
		Host: "test.backstage.dev",
		Port: ":4567",
	}

	service = &account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage"}
}

var _ = Suite(&S{})
