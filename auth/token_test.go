package auth

import (
	"github.com/albertoleal/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestGenerateAndGetToken(c *C) {
	user := &account.User{Username: "alice"}
	tokenInfo := GenerateToken(user)
	tokenType, token, err := GetToken("Token " + tokenInfo.Token)
	c.Assert(tokenType, Equals, "Token")
	c.Assert(token, Equals, tokenInfo.Token)
	c.Assert(err, IsNil)
}

func (s *S) TestGenerateAndGetTokenWithInvalidTokenType(c *C) {
	user := &account.User{Username: "alice"}
	tokenInfo := GenerateToken(user)
	tokenType, token, err := GetToken("InvalidType " + tokenInfo.Token)
	c.Assert(tokenType, Equals, "")
	c.Assert(token, Equals, "")
	c.Assert(err.Error(), Equals, "Invalid token format.")
}

func (s *S) TestGetTokenWithInvalidFormat(c *C) {
	tokenType, token, err := GetToken("Invalid-Format")
	c.Assert(tokenType, Equals, "")
	c.Assert(token, Equals, "")
	c.Assert(err.Error(), Equals, "Invalid token format.")
}

func (s *S) TestGenerateToken(c *C) {
	token := GenerateToken(&account.User{Username: "alice"})
	c.Assert(len(token.Token), Equals, 44)
	c.Assert(token.Type, Equals, "Token")
}
