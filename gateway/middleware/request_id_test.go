package middleware

import (
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *S) TestRequestId(c *C) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(len(r.Header.Get("X-Request-Id")), Not(Equals), 0)
	})

	rid := &RequestIdMiddleware{}
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://apihub.example.org", nil)
	rid.ServeHTTP(res, req, next)
}
