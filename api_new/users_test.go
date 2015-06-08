package api_new_test

import (
	"net/http"

	"github.com/backstage/backstage/api_new"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	headers, code, body, err := httpClient.MakeRequest(api_new.RequestArgs{
		Method: "POST",
		Path:   "/auth/signup",
		Body:   `{"name": "Alice", "email": "alice@example.org", "password": "123456"}`,
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"Alice","email":"alice@example.org"}`)
}

func (s *S) TestCreateUserWithInvalidPayloadFormat(c *C) {
	headers, code, body, err := httpClient.MakeRequest(api_new.RequestArgs{
		Method: "POST",
		Path:   "/auth/signup",
		Body:   `"name": "Alice"`,
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestCreateUserWithMissingRequiredFields(c *C) {
	headers, code, body, err := httpClient.MakeRequest(api_new.RequestArgs{
		Method: "POST",
		Path:   "/auth/signup",
		Body:   `{}`,
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Name/Email/Password cannot be empty."}`)
}
