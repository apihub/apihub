package api_test

import (
	"net/http"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestAddGetRequestError(c *C) {
	req, _ := http.NewRequest("GET", "/foo", nil)
	api.AddRequestError(req, errors.ErrClientNotFound)
	err, ok := api.GetRequestError(req)
	c.Assert(err, Equals, errors.ErrClientNotFound)
	c.Assert(ok, Equals, true)
}

func (s *S) TestSetAndGetCurrentUser(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	req, _ := http.NewRequest("GET", "/foo", nil)
	api.SetCurrentUser(req, user)
	u, err := api.GetCurrentUser(req)
	c.Assert(u, DeepEquals, user)
	c.Check(err, IsNil)
}

func (s *S) TestGetCurrentUserNotSignedIn(c *C) {
	req, _ := http.NewRequest("GET", "/foo", nil)
	u, err := api.GetCurrentUser(req)
	c.Check(u, IsNil)
	c.Assert(err, Equals, errors.ErrLoginRequired)
}
