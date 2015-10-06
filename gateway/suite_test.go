package gateway

import (
	"testing"

	"github.com/apihub/apihub/account"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	Settings *Settings
}

var service *account.Service

func (s *S) SetUpTest(c *C) {
	s.Settings = &Settings{
		Host: "test.apihub.dev",
		Port: ":4567",
	}

	service = &account.Service{Endpoint: "http://example.org/api", Subdomain: "apihub"}
}

var _ = Suite(&S{})
