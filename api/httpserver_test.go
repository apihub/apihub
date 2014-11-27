package api

import (
	"net/http"

	. "gopkg.in/check.v1"
)

func (s *S) TestAuthorizationMiddlewareEnabled(c *C) {
	req, err := http.NewRequest("POST", "/services", nil)
	if err != nil {
		c.Error(err)
	}

	s.server.mux.ServeHTTP(s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 401)
}

func (s *S) TestAuthorizationMiddlewareDisabledUnderDebug(c *C) {
	req, err := http.NewRequest("GET", "/debug/helloworld", nil)
	if err != nil {
		c.Error(err)
	}

	s.server.mux.ServeHTTP(s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 200)
}
