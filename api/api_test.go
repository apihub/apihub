package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

type HiHandler struct {
	ApiHandler
}

func (handler *HiHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return OK("Hi from custom route!")
}

func (s *S) TestAddPrivateRoute(c *C) {
	s.Api.init()
	s.Api.AddPrivateRoute("GET", "/hi", &HiHandler{}, "Index")
	req, _ := http.NewRequest("GET", "/hi", nil)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.Api.privateRoutes.ServeHTTPC(webC, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, "Hi from custom route!")
}
