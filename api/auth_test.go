package api

import (
	. "github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestLogin(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	u := &User{Email: "alice@example.org", Password: "123456"}
	_, ok := Login(u)
	c.Assert(ok, IsNil)
}

func (s *S) TestLoginWithInvalidUsername(c *C) {
	user := &User{Username: "invalid-email", Password: "123456"}
	_, ok := Login(user)
	c.Assert(ok, NotNil)
}

func (s *S) TestLoginWithInvalidPassword(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()
	_, ok := Login(user)
	c.Assert(ok, NotNil)
}

func (s *S) TestLoginAndGetToken(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	u := &User{Email: "alice@example.org", Password: "123456"}
	token, ok := LoginAndGetToken(u)
	c.Assert(ok, IsNil)
	c.Assert(token, NotNil)
}

func (s *S) TestLoginAndGetTokenWithInvalidCredentials(c *C) {
	user := &User{Username: "invalid-email", Password: "123456"}
	token, ok := LoginAndGetToken(user)
	c.Assert(ok, NotNil)
	c.Assert(token, IsNil)
}
