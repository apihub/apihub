package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/backstage/backstage/auth"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestAuthorizationMiddlewareWithValidToken(c *C) {
	s.router.Use(ErrorMiddleware)
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
	c.Assert(s.recorder.Body.String(), Equals, "")
}

func (s *S) TestAuthorizationMiddlewareWithInvalidToken(c *C) {
	s.router.Use(ErrorMiddleware)
	s.router.Use(AuthorizationMiddleware)
	s.router.Get("/", s.handler)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Invalid-Token")
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	erro, _ := GetRequestError(&cc)
	c.Assert(erro.StatusCode, Equals, http.StatusUnauthorized)
	c.Assert(erro.ErrorDescription, Equals, "Request refused or access is not allowed.")
	c.Assert(erro.ErrorType, Equals, "unauthorized_access")
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"unauthorized_access","error_description":"Request refused or access is not allowed."}`)
}

func (s *S) TestAuthorizationMiddlewareWithMissingToken(c *C) {
	s.router.Use(ErrorMiddleware)
	s.router.Use(AuthorizationMiddleware)
	s.router.Get("/", s.handler)
	req, _ := http.NewRequest("GET", "/", nil)
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	erro, _ := GetRequestError(&cc)
	c.Assert(erro.StatusCode, Equals, http.StatusUnauthorized)
	c.Assert(erro.ErrorDescription, Equals, "Request refused or access is not allowed.")
	c.Assert(erro.ErrorType, Equals, "unauthorized_access")
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"unauthorized_access","error_description":"Request refused or access is not allowed."}`)
}

func (s *S) TestRequestIdMiddleware(c *C) {
	s.router.Use(ErrorMiddleware)
	s.router.Use(RequestIdMiddleware)
	s.router.Get("/", s.handler)

	req, _ := http.NewRequest("GET", "/", nil)
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, http.StatusOK)
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
	body := &HTTPResponse{}
	json.Unmarshal(w.Body.Bytes(), body)
	c.Assert(body.ErrorDescription, Equals, "The resource requested does not exist.")
	c.Assert(body.ErrorType, Equals, "not_found")
}
