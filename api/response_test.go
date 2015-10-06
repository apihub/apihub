package api_test

import (
	"net/http"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestToJsonWithError(c *C) {
	erro := errors.ErrorResponse{
		Type:        "invalid_request",
		Description: "The request is missing a required parameter.",
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
