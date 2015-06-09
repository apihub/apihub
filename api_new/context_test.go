package api_new_test

import (
	"net/http"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/api_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestAddGetRequestError(c *C) {
	req, _ := http.NewRequest("GET", "/foo", nil)
	api_new.AddRequestError(req, errors.ErrClientNotFound)
	err, ok := api_new.GetRequestError(req)
	c.Assert(err, Equals, errors.ErrClientNotFound)
	c.Assert(ok, Equals, true)
}

func (s *S) TestSetAndGetCurrentUser(c *C) {
	user := &account_new.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	req, _ := http.NewRequest("GET", "/foo", nil)
	api_new.SetCurrentUser(req, user)
	u, err := api_new.GetCurrentUser(req)
	c.Assert(u, DeepEquals, user)
	c.Check(err, IsNil)
}

func (s *S) TestGetCurrentUserNotSignedIn(c *C) {
	req, _ := http.NewRequest("GET", "/foo", nil)
	u, err := api_new.GetCurrentUser(req)
	c.Check(u, IsNil)
	c.Assert(err, Equals, errors.ErrLoginRequired)
}
