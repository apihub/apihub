package gateway

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestGatewayNotFound(c *C) {
	gateway := NewGateway(s.config, nil)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "invalid.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusNotFound)
	c.Assert(w.Body.String(), Equals, "Not found.\n")
}

func (s *S) TestGatewayExistingService(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}}
	gateway := NewGateway(s.config, services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
}
