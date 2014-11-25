package app

import (
	. "gopkg.in/check.v1"

	"github.com/albertoleal/backstage/errors"
)

func (s *S) TestCreateUser(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	err := CreateUser(&user)
	defer DeleteUser(&user)
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := User{}
	err := CreateUser(&user)

	e := err.(*errors.ValidationError)
	msg := "Name/Email/Username/Password cannot be empty."
	c.Assert(e.Message, Equals, msg)
}
