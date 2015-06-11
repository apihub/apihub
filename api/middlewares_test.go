package api_test

import (
	"fmt"
	"net/http"

	. "gopkg.in/check.v1"
)

func (s *S) TestAuthorizationMiddleware(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/users",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"%s","email":"%s"}`, user.Name, user.Email))
}

func (s *S) TestAuthorizationMiddlewareWithInvalidToken(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/users", Headers: http.Header{"Authorization": {"expired-token"}}}, c)
}

func (s *S) TestAuthorizationMiddlewareWithMissingToken(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/users"}, c)
}

func (s *S) TestNotFoundHandler(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "GET",
		Path:   "/not-found-path",
	})

	c.Check(err, IsNil)
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"The resource requested does not exist."}`)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *S) TestRequestId(c *C) {
	headers, _, _, _ := httpClient.MakeRequest(RequestArgs{
		Method: "DELETE",
		Path:   "/api/users",
	})

	c.Assert(headers.Get("X-Request-Id"), Not(Equals), "")
}
