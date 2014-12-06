package auth

import (
	"github.com/albertoleal/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestGenerateAndGetUserFromToken(c *C) {
	user := &account.User{Username: "alice", Email: "alice@bar.example.org"}
	tokenInfo := GenerateToken(user)
	u, err := GetUserFromToken("Token " + tokenInfo.Token)
	c.Assert(err, IsNil)
	c.Assert(u.Email, Equals, "alice@bar.example.org")
}

func (s *S) TestGenerateAndGetUsernameFromTokenWithInvalidTokenType(c *C) {
	user := &account.User{Username: "alice", Email: "alice@bar.example.org"}
	tokenInfo := GenerateToken(user)
	user, err := GetUserFromToken("InvalidType " + tokenInfo.Token)
	c.Assert(err.Error(), Equals, "Invalid token format.")
	c.Assert(user, IsNil)
}

func (s *S) TestGetUsernameFromTokenWithInvalidFormat(c *C) {
	user, err := GetUserFromToken("Invalid-Format")
	c.Assert(err.Error(), Equals, "Invalid token format.")
	c.Assert(user, IsNil)
}

func (s *S) TestGenerateToken(c *C) {
	token := GenerateToken(&account.User{Email: "alice@bar.example.org"})
	c.Assert(token.Type, Equals, "Token")
	c.Assert(len(token.Token), Equals, 44)
}
