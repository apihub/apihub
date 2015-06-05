package auth

import (
	"github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

var user = &account.User{Username: "alice", Email: "alice@bar.example.org"}

func (s *S) TestGenerateAndGetUserFromToken(c *C) {
	tokenInfo := GenerateToken(user)
	u, err := GetUserFromToken("Token " + tokenInfo.Token)
	c.Assert(err, IsNil)
	c.Assert(u.Email, Equals, "alice@bar.example.org")
}

func (s *S) TestGenerateAndGetUserFromTokenWithInvalidTokenType(c *C) {
	tokenInfo := GenerateToken(user)
	user, err := GetUserFromToken("InvalidType " + tokenInfo.Token)
	c.Assert(err.Error(), Equals, "Invalid token format.")
	c.Assert(user, IsNil)
}

func (s *S) TestGetUserFromTokenWithInvalidFormat(c *C) {
	user, err := GetUserFromToken("Invalid-Format")
	c.Assert(err.Error(), Equals, "Invalid token format.")
	c.Assert(user, IsNil)
}

func (s *S) TestGetUserFromTokenWithNonExistingToken(c *C) {
	user, err := GetUserFromToken("Token xyz")
	c.Assert(err.Error(), Equals, "Token not found.")
	c.Assert(user, IsNil)
}

func (s *S) TestGenerateToken(c *C) {
	token := GenerateToken(user)
	c.Assert(token.Type, Equals, "Token")
	c.Assert(len(token.Token), Equals, 44)
}

func (s *S) TestTokenForReturnsTheSameTokenIfValid(c *C) {
	tokenInfo := TokenFor(user)
	tokenInfo2 := TokenFor(user)
	c.Assert(tokenInfo.Token, Equals, tokenInfo2.Token)
}

func (s *S) TestRevokeTokensFor(c *C) {
	tokenInfo := TokenFor(user)
	_, err := GetUserFromToken(tokenInfo.Type + " " + tokenInfo.Token)
	c.Assert(err, IsNil)
	RevokeTokensFor(user)
	_, err = GetUserFromToken(tokenInfo.Type + " " + tokenInfo.Token)
	c.Assert(err, NotNil)
}

func (s *S) TestTokenToString(c *C) {
	tokenInfo := GenerateToken(user)
	str := tokenInfo.ToString()
	c.Assert(str, Matches, `^{"access_token":".*?","token_type":"Token","expires":86400,"created_at":".*?"}$`)
}
