package api_new_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/backstage/backstage/api_new"
	. "gopkg.in/check.v1"
)

func (s *S) TestMakeRequest(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Alice"}`))
	}))
	defer server.Close()
	httpClient.Host = server.URL

	args := api_new.RequestArgs{Method: "GET", Path: "/path", Body: ""}
	_, _, body, err := httpClient.MakeRequest(args)
	c.Assert(string(body), Equals, `{"name": "Alice"}`)
	c.Check(err, IsNil)
}

func (s *S) TestReturnsErrorWhenHostIsInvalid(c *C) {
	httpClient.Host = "://invalid-host"
	args := api_new.RequestArgs{Method: "GET", Path: "/path", Body: ""}
	_, _, _, err := httpClient.MakeRequest(args)
	c.Check(err, Not(IsNil))
}

func (s *S) TestReturnsErrorWhenRequestIsInvalid(c *C) {
	httpClient.Host = "invalid-host"
	args := api_new.RequestArgs{Method: "GET", Path: "/path", Body: ""}
	_, _, _, err := httpClient.MakeRequest(args)
	c.Check(err, Not(IsNil))
}

func (s *S) TestReturnsErrorWhenResponseIsInvalid(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Length", "1")
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	httpClient.Host = server.URL

	args := api_new.RequestArgs{Method: "GET", Path: "/path", Body: ""}
	_, _, _, err := httpClient.MakeRequest(args)
	c.Check(err, Not(IsNil))
}
