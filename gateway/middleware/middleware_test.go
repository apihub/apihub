package middleware

import (
	. "gopkg.in/check.v1"
)

func (s *S) TestGetMiddleware(c *C) {
	c.Check(s.middlewares.Get("invalid"), IsNil)
}

func (s *S) TestAddMiddleware(c *C) {
	c.Check(s.middlewares.Get("CustomMiddleware"), IsNil)
	s.middlewares.Add("CustomMiddleware", NewCorsMiddleware)
	c.Check(s.middlewares.Get("CustomMiddleware"), NotNil)
}
