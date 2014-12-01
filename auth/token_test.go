package auth

import (
	. "gopkg.in/check.v1"
)

func (s *S) TestGetToken(c *C) {
	tokenType, token, err := GetToken("Basic token")
	c.Assert(tokenType, Equals, "Basic")
	c.Assert(token, Equals, "token")
	c.Assert(err, IsNil)
}

func (s *S) TestGetTokenWithInvalidFormat(c *C) {
	tokenType, token, err := GetToken("Invalid-Format")
	c.Assert(tokenType, Equals, "")
	c.Assert(token, Equals, "")
	c.Assert(err.Error(), Equals, "Invalid token format.")
}

func (s *S) TestGenerateToken(c *C) {
	token := GenerateToken()
	c.Assert(len(token.Token), Equals, 44)
}
