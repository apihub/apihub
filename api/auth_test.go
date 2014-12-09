package api

import (
	. "github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestLogin(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()
	_, ok := Login(user.Email, "123456")
	c.Assert(ok, IsNil)
}

func (s *S) TestLoginWithInvalidUsername(c *C) {
	_, ok := Login("invalid-email", "123456")
	c.Assert(ok, NotNil)
}

func (s *S) TestLoginWithInvalidPassword(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()
	_, ok := Login(user.Email, "invalid-password")
	c.Assert(ok, NotNil)
}
