package api

import (
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *S) SetUpTest(c *C) {
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	s.recorder = httptest.NewRecorder()
}

func (s *S) TestAuthorizationMiddlewareWithValidToken(c *C) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		c.Error(err)
	}

	req.Header.Set("Authorization", "Basic xyz")
	authorizationMiddleware(s.recorder, req, s.handler)
	c.Assert(s.recorder.Code, Equals, 200)
}

func (s *S) TestAuthorizationMiddlewareWithInvalidToken(c *C) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		c.Error(err)
	}

	req.Header.Set("Authorization", "Invalid-Token")
	authorizationMiddleware(s.recorder, req, s.handler)
	c.Assert(s.recorder.Code, Equals, 401)
}

func (s *S) TestAuthorizationMiddlewareWithMissingToken(c *C) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		c.Error(err)
	}

	authorizationMiddleware(s.recorder, req, s.handler)
	c.Assert(s.recorder.Code, Equals, 401)
}
