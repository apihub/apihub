package middleware

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	middlewares Middlewares
}

func (s *S) SetUpTest(c *C) {
		s.middlewares =  map[string]func() Middleware{}
}

var _ = Suite(&S{})
