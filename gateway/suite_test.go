package gateway

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	settings *Settings
}

func (s *S) SetUpTest(c *C) {
	s.settings = &Settings{
		Host:        "test.backstage.dev",
		Port:        ":4567",
		ChannelName: "services",
	}
}

var _ = Suite(&S{})
