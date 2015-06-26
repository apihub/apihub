package gateway

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/backstage/maestro/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestGatewayNotFound(c *C) {
	gateway := New(s.Settings, nil)
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "invalid.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusNotFound)
	c.Assert(w.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(w.Body.String(), Equals, "{\"error\":\"not_found\",\"error_description\":\"The requested resource could not be found but may be available again in the future.\"}\n")
}

func (s *S) TestGatewayExistingService(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		c.Assert(len(r.Header.Get("X-Request-Id")), Not(Equals), 0)
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}}
	gateway := New(s.Settings, nil)
	gateway.LoadServices(services)
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Body.String(), Equals, "OK")
	c.Assert(w.Code, Equals, http.StatusOK)
}

func (s *S) TestGatewayTimeout(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Timeout: 1}}
	gateway := New(s.Settings, nil)
	gateway.LoadServices(services)

	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusGatewayTimeout)
	c.Assert(w.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(w.Body.String(), Equals, `{"error":"gateway_timeout","error_description":"The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server."}`)
}

func (s *S) TestGatewayCopyResponseHeaders(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}}
	gateway := New(s.Settings, nil)
	gateway.LoadServices(services)

	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusCreated)
	c.Assert(w.Header().Get("Content-Type"), Equals, "text/plain; charset=utf-8")
}

func (s *S) TestGatewayInternalError(c *C) {
	services := []*account.Service{&account.Service{Endpoint: "http://invalidurl", Subdomain: "test"}}
	gateway := New(s.Settings, nil)
	gateway.LoadServices(services)

	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
	c.Assert(w.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(w.Body.String(), Equals, `{"error":"internal_server_error","error_description":"dial tcp: lookup invalidurl: no such host"}`)
}

func (s *S) TestAddService(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		c.Assert(len(r.Header.Get("X-Request-Id")), Not(Equals), 0)
	}))
	defer target.Close()

	service := &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}
	gateway := New(s.Settings, nil)
	gateway.AddService(service)
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Body.String(), Equals, "OK")
	c.Assert(w.Code, Equals, http.StatusOK)
}

func (s *S) TestRemoveService(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		c.Assert(len(r.Header.Get("X-Request-Id")), Not(Equals), 0)
	}))
	defer target.Close()

	service := &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}
	gateway := New(s.Settings, nil)
	gateway.AddService(service)
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)
	c.Assert(w.Body.String(), Equals, "OK")
	c.Assert(w.Code, Equals, http.StatusOK)

	gateway.RemoveService(service)
	r, _ = http.NewRequest("GET", "http://test.backstage.dev", nil)
	w = httptest.NewRecorder()
	gateway.ServeHTTP(w, r)
	c.Assert(w.Body.String(), Equals, "{\"error\":\"not_found\",\"error_description\":\"The requested resource could not be found but may be available again in the future.\"}\n")
	c.Assert(w.Code, Equals, http.StatusNotFound)
}
