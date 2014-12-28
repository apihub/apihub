package api

import (
	"bytes"
	"io/ioutil"

	"github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestGetCurrentUserWithoutUser(c *C) {
	env := map[string]interface{}{}
	webContext := &web.C{Env: env}
	api := &ApiHandler{}
	user, err := api.getCurrentUser(webContext)
	c.Assert(err, NotNil)
	c.Assert(user, IsNil)
}

func (s *S) TestGetCurrentUserWhenUserDoesNotExist(c *C) {
	alice := &account.User{Username: "alice", Name: "Alice", Email: "alice@example.org", Password: "123456"}
	env := map[string]interface{}{CurrentUser: alice}
	webContext := &web.C{Env: env}
	api := &ApiHandler{}
	user, err := api.getCurrentUser(webContext)
	c.Assert(err, NotNil)
	c.Assert(user, IsNil)
}

func (s *S) TestGetCurrentUser(c *C) {
	alice := &account.User{Username: "alice", Name: "Alice", Email: "alice@example.org", Password: "123456"}
	alice.Save()
	defer alice.Delete()

	env := map[string]interface{}{CurrentUser: alice}
	webContext := &web.C{Env: env}
	api := &ApiHandler{}
	user, err := api.getCurrentUser(webContext)
	c.Assert(user.Name, Equals, "Alice")
	c.Assert(err, IsNil)
}

func (s *S) TestParseBody(c *C) {
	api := &ApiHandler{}
	var result map[string]string
	body := ioutil.NopCloser(bytes.NewBufferString(`{"name": "Alice"}`))
	api.parseBody(body, &result)
	c.Assert(result["name"], Equals, "Alice")
}

func (s *S) TestParseBodyWithInvalidInterface(c *C) {
	api := &ApiHandler{}
	var result map[string]string
	body := ioutil.NopCloser(bytes.NewBufferString(`"name": "Alice"`))
	err := api.parseBody(body, &result)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "The request was invalid or cannot be served.")
}
