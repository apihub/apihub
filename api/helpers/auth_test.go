package helpers

import (
	"testing"

	. "github.com/albertoleal/backstage/account"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestSignIn(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer DeleteUser(user)
	_, ok := SignIn(user.Username, "123456")
	c.Assert(ok, IsNil)
}

func (s *S) TestSignInWithInvalidUsername(c *C) {
	_, ok := SignIn("invalid-username", "123456")
	c.Assert(ok, NotNil)
}

func (s *S) TestSignInWithInvalidPassword(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer DeleteUser(user)
	_, ok := SignIn(user.Username, "invalid-password")
	c.Assert(ok, NotNil)
}
