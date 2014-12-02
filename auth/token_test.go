package auth

import (
	"github.com/albertoleal/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestGenerateAndGetUserFromToken(c *C) {
	user := &account.User{Username: "alice"}
	tokenInfo := GenerateToken(user)
	u, err := GetUserFromToken("Token " + tokenInfo.Token)
	c.Assert(err, IsNil)
	c.Assert(u.Username, Equals, "alice")
}

func (s *S) TestGenerateAndGetUsernameFromTokenWithInvalidTokenType(c *C) {
	user := &account.User{Username: "alice"}
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
	token := GenerateToken(&account.User{Username: "alice"})
	c.Assert(token.Type, Equals, "Token")
	c.Assert(len(token.Token), Equals, 44)
}
