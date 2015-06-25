package middleware

import (
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *S) TestCors(c *C) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	cors := &Cors{}
	cors.AllowedOrigins = []string{"http://backstage.example.org"}
	cors.AllowedMethods = []string{"GET", "POST"}
	cors.AllowedHeaders = []string{"origin"}
	cors.ExposedHeaders = []string{}
	cors.AllowCredentials = true
	cors.MaxAge = 100
	cors.Debug = true

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://backstage.example.org", nil)
	req.Header.Add("Origin", "http://backstage.example.org")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")

	cors.ProcessRequest(res, req, next)
	c.Assert(res.Header().Get("Access-Control-Allow-Origin"), Equals, "http://backstage.example.org")
	c.Assert(res.Header().Get("Access-Control-Allow-Methods"), Equals, "GET")
	c.Assert(res.Header().Get("Access-Control-Allow-Headers"), Equals, "Origin")
	c.Assert(res.Header().Get("Access-Control-Allow-Credentials"), Equals, "true")
	c.Assert(res.Header().Get("Access-Control-Max-Age"), Equals, "100")
	c.Assert(res.Header().Get("Vary"), Equals, "Origin")
}
