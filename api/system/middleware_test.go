package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/albertoleal/backstage/api/context"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) SetUpTest(c *C) {
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	s.recorder = httptest.NewRecorder()
	s.router = web.New()
	s.router.Use(AuthorizationMiddleware)
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

func (s *S) TestAuthorizationMiddlewareWithValidToken(c *C) {
	s.router.Get("/", s.handler)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer xyz")
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	_, ok := context.GetRequestError(&cc)
	c.Assert(ok, Equals, false)
}

func (s *S) TestAuthorizationMiddlewareWithInvalidToken(c *C) {
	s.router.Get("/", s.handler)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Invalid-Token")
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	erro, _ := context.GetRequestError(&cc)
	c.Assert(erro.StatusCode, Equals, http.StatusUnauthorized)
	c.Assert(erro.Message, Equals, "You do not have access to this resource.")
}

func (s *S) TestAuthorizationMiddlewareWithMissingToken(c *C) {
	s.router.Get("/", s.handler)

	req, _ := http.NewRequest("GET", "/", nil)
	cc := web.C{Env: map[string]interface{}{}}
	s.router.ServeHTTPC(cc, s.recorder, req)
	erro, _ := context.GetRequestError(&cc)
	c.Assert(erro.StatusCode, Equals, http.StatusUnauthorized)
	c.Assert(erro.Message, Equals, "You do not have access to this resource.")
}
