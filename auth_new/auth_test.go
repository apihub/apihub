package auth_new_test

import (
	"fmt"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestAuthenticate(c *C) {
	user := &account_new.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	found, ok := s.auth.Authenticate("alice@bar.example.org", "123")
	c.Assert(found, DeepEquals, user)
	c.Assert(ok, Equals, true)
}

func (s *S) TestAuthenticateWithInvalidCredentials(c *C) {
	user := account_new.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	_, ok := s.auth.Authenticate(user.Email, "invalid-password")
	c.Assert(ok, Equals, false)
}

func (s *S) TestAuthenticateWithNotFound(c *C) {
	_, ok := s.auth.Authenticate("invalid-email", "invalid-password")
	c.Assert(ok, Equals, false)
}

func (s *S) TestUserFromToken(c *C) {
	user := &account_new.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	token, _ := s.auth.Login("alice@bar.example.org", "123")
	foundUser, err := s.auth.UserFromToken(fmt.Sprintf("%s %s", token.Type, token.Token))
	c.Check(err, IsNil)
	c.Assert(foundUser, DeepEquals, user)
}

func (s *S) TestUserFromTokenWithNotFound(c *C) {
	foundUser, err := s.auth.UserFromToken("Token invalid-token")
	c.Check(err, Not(IsNil))
	c.Check(foundUser, IsNil)
}

func (s *S) TestUserFromTokenWithInvalidFormat(c *C) {
	foundUser, err := s.auth.UserFromToken("invalid-format")
	c.Check(err, Not(IsNil))
	c.Check(foundUser, IsNil)
}

func (s *S) TestRevokeUserToken(c *C) {
	user := &account_new.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	token, _ := s.auth.Login("alice@bar.example.org", "123")
	tstr := fmt.Sprintf("%s %s", token.Type, token.Token)
	foundUser, err := s.auth.UserFromToken(tstr)
	c.Check(err, IsNil)
	c.Assert(foundUser, DeepEquals, user)

	s.auth.RevokeUserToken(tstr)
	foundUser, err = s.auth.UserFromToken(tstr)
	c.Assert(err, Equals, errors.ErrTokenNotFound)
	c.Check(foundUser, IsNil)
}
