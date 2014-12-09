package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/backstage/backstage/auth"
	"github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestAuthorizationMiddlewareWithValidToken(c *C) {
	s.router.Use(AuthorizationMiddleware)
	err := alice.Save()
	defer alice.Delete()
	if err != nil {
		c.Error(err)
	}

	tokenInfo := auth.GenerateToken(alice)
	s.router.Get("/", s.handler)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Token "+tokenInfo.Token)
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	_, ok := GetRequestError(&cc)
	c.Assert(ok, Equals, false)
}

func (s *S) TestAuthorizationMiddlewareWithInvalidToken(c *C) {
	s.router.Use(AuthorizationMiddleware)
	s.router.Get("/", s.handler)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Invalid-Token")
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	erro, _ := GetRequestError(&cc)
	c.Assert(erro.StatusCode, Equals, http.StatusUnauthorized)
	c.Assert(erro.Payload, Equals, "You do not have access to this resource.")
}

func (s *S) TestAuthorizationMiddlewareWithMissingToken(c *C) {
	s.router.Use(AuthorizationMiddleware)
	s.router.Get("/", s.handler)
	req, _ := http.NewRequest("GET", "/", nil)
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	erro, _ := GetRequestError(&cc)
	c.Assert(erro.StatusCode, Equals, http.StatusUnauthorized)
	c.Assert(erro.Payload, Equals, "You do not have access to this resource.")
}

func (s *S) TestRequestIdMiddleware(c *C) {
	s.router.Use(RequestIdMiddleware)
	s.router.Get("/", s.handler)

	req, _ := http.NewRequest("GET", "/", nil)
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.HeaderMap["Request-Id"], NotNil)
}

func (s *S) TestNotFoundHandler(c *C) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/invalid-endpoint", nil)
	if err != nil {
		c.Error(err)
	}

	NotFoundHandler(w, req)
	c.Assert(w.Code, Equals, http.StatusNotFound)
	body := &errors.HTTPError{}
	json.Unmarshal(w.Body.Bytes(), body)
	c.Assert(body.Message, Equals, "The resource you are looking for was not found.")
}
