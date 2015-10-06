package test

import (
	"fmt"
	"testing"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/auth"
	"github.com/apihub/apihub/errors"
	. "gopkg.in/check.v1"
)

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type AuthenticatableSuite struct {
	Auth auth.Authenticatable
}

func (s *AuthenticatableSuite) TestAuthenticate(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	found, ok := s.Auth.Authenticate("alice@bar.example.org", "123")
	c.Assert(found, DeepEquals, user)
	c.Assert(ok, Equals, true)
}

func (s *AuthenticatableSuite) TestAuthenticateWithInvalidCredentials(c *C) {
	user := account.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	_, ok := s.Auth.Authenticate(user.Email, "invalid-password")
	c.Assert(ok, Equals, false)
}

func (s *AuthenticatableSuite) TestAuthenticateWithNotFound(c *C) {
	_, ok := s.Auth.Authenticate("invalid-email", "invalid-password")
	c.Assert(ok, Equals, false)
}

func (s *AuthenticatableSuite) TestCreateUserToken(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	token, err := s.Auth.CreateUserToken(user)
	c.Check(err, IsNil)
	c.Assert(token.AccessToken, Not(Equals), "")
}

func (s *AuthenticatableSuite) TestUserFromToken(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	token, _ := s.Auth.CreateUserToken(user)
	foundUser, err := s.Auth.UserFromToken(fmt.Sprintf("%s %s", token.Type, token.AccessToken))
	c.Check(err, IsNil)
	c.Assert(foundUser, DeepEquals, user)
}

func (s *AuthenticatableSuite) TestUserFromTokenWithNotFound(c *C) {
	foundUser, err := s.Auth.UserFromToken("Token invalid-token")
	c.Check(err, Not(IsNil))
	c.Check(foundUser, IsNil)
}

func (s *AuthenticatableSuite) TestUserFromTokenWithInvalidFormat(c *C) {
	foundUser, err := s.Auth.UserFromToken("invalid-format")
	c.Check(err, Not(IsNil))
	c.Check(foundUser, IsNil)
}

func (s *AuthenticatableSuite) TestRevokeUserToken(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@bar.example.org", Password: "123"}
	user.Create()
	defer user.Delete()

	token, _ := s.Auth.CreateUserToken(user)
	tstr := fmt.Sprintf("%s %s", token.Type, token.AccessToken)
	foundUser, err := s.Auth.UserFromToken(tstr)
	c.Check(err, IsNil)
	c.Assert(foundUser, DeepEquals, user)

	s.Auth.RevokeUserToken(tstr)
	foundUser, err = s.Auth.UserFromToken(tstr)
	c.Assert(err, Equals, errors.ErrTokenNotFound)
	c.Check(foundUser, IsNil)
}
