package middleware

import (
	. "gopkg.in/check.v1"
	"net/http"
)

func (s *S) TestGetMiddleware(c *C) {
	c.Check(s.middlewares.Get("invalid"), IsNil)
}

func (s *S) TestAddMiddleware(c *C) {
	c.Check(s.middlewares.Get("CustomMiddleware"), IsNil)
	ah := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}
	s.middlewares.Add("CustomMiddleware", ah)
	c.Check(s.middlewares.Get("CustomMiddleware"), NotNil)
}
