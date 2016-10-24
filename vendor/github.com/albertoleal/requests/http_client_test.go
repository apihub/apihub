package requests_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/albertoleal/requests"
	. "gopkg.in/check.v1"
)

func (s *S) TestMakeRequest(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Alice"}`))
	}))
	defer server.Close()
	httpClient := requests.NewHTTPClient(server.URL)

	args := requests.Args{Method: "GET", Path: "/path", Body: nil, AcceptableCode: http.StatusOK}
	_, _, body, err := httpClient.MakeRequest(args)
	c.Assert(string(body), Equals, `{"name": "Alice"}`)
	c.Check(err, IsNil)
}

func (s *S) TestMakeRequestWithNonAcceptableCode(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Alice"}`))
	}))
	defer server.Close()
	httpClient := requests.NewHTTPClient(server.URL)

	args := requests.Args{Method: "GET", Path: "/path", Body: nil, AcceptableCode: http.StatusBadRequest}
	_, _, body, err := httpClient.MakeRequest(args)
	c.Assert(string(body), Equals, `{"name": "Alice"}`)
	e, ok := err.(requests.ResponseError)
	c.Assert(ok, Equals, true)
	c.Assert(e.Error(), Equals, "The response was invalid or cannot be served.")
}

func (s *S) TestReturnsErrorWhenHostIsInvalid(c *C) {
	httpClient := requests.NewHTTPClient("://invalid-host")
	args := requests.Args{Method: "GET", Path: "/path", Body: nil}
	_, _, _, err := httpClient.MakeRequest(args)
	_, ok := err.(requests.InvalidHostError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsErrorWhenRequestIsInvalid(c *C) {
	httpClient := requests.NewHTTPClient("invalid-host")
	args := requests.Args{Method: "GET", Path: "/path", Body: nil}
	_, _, _, err := httpClient.MakeRequest(args)
	_, ok := err.(requests.RequestError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestReturnsErrorWhenResponseIsInvalid(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Length", "1")
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	httpClient := requests.NewHTTPClient(server.URL)

	args := requests.Args{Method: "GET", Path: "/path", Body: nil}
	_, _, _, err := httpClient.MakeRequest(args)
	_, ok := err.(requests.ResponseError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestIncludesHeader(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		c.Assert(auth, Equals, "Token abcde")
	}))
	defer server.Close()

	httpClient := requests.NewHTTPClient(server.URL)

	args := requests.Args{Method: "GET", Path: "/path", Body: nil, Headers: http.Header{"Authorization": {"Token abcde"}}}
	httpClient.MakeRequest(args)
}

func (s *S) TestReturnsDefaultError(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	httpClient := requests.NewHTTPClient(server.URL)

	args := requests.Args{Method: "GET", Path: "/path", Body: nil}
	_, _, _, err := httpClient.MakeRequest(args)
	e, ok := err.(requests.ResponseError)
	c.Assert(e.Error(), Equals, "The response was invalid or cannot be served.")
	c.Assert(ok, Equals, true)
}
