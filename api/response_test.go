package api_test

import (
	"net/http"

	"github.com/backstage/apimanager/account"
	"github.com/backstage/apimanager/api"
	. "gopkg.in/check.v1"
)

func (s *S) TestToJsonWithError(c *C) {
	erro := api.HTTPError{
		ErrorType:        "invalid_request",
		ErrorDescription: "The request is missing a required parameter.",
	}

	err := &api.HTTPResponse{
		StatusCode: http.StatusBadRequest,
		Body:       erro,
	}
	c.Assert(string(err.ToJson()), Equals, `{"error":"invalid_request","error_description":"The request is missing a required parameter."}`)
}

func (s *S) TestToJsonWithUser(c *C) {
	alice := &account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	err := &api.HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       alice,
	}
	c.Assert(string(err.ToJson()), Equals, `{"name":"Alice","email":"alice@example.org","password":"123456"}`)
}
