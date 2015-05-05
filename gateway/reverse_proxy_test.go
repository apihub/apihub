package gateway

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *S) TestServer(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	service := &Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Timeout: 10}
	rp := NewReverseProxy(service)

	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "test.backstage.dev", nil)
	rp.proxy.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
}
