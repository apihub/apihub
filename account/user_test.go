package account

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

func (s *S) TestCreateUserWithSameUsername(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	CreateUser(&user)
	defer DeleteUser(&user)

	user2 := User{Name: "Bob", Email: "bob@bar.com", Username: "alice", Password: "123456"}
	err := CreateUser(&user2)
	e := err.(*errors.ValidationError)
	msg := "Someone already has that username. Could you try another?."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := User{}
	err := CreateUser(&user)

	e := err.(*errors.ValidationError)
	msg := "Name/Email/Username/Password cannot be empty."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCreateUserShouldMaskThePassword(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	CreateUser(&user)
	defer DeleteUser(&user)

	foundUser, _ := FindUserByUsername("alice")
	c.Assert(foundUser.Password, Not(Equals), "123456")
}

func (s *S) FindUserByUsername(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	CreateUser(&user)
	defer DeleteUser(&user)

	foundUser, err := FindUserByUsername("alice")
	c.Assert(err, IsNil)
	c.Assert(foundUser, NotNil)
}

func (s *S) TestFindUserWithInvalidUsername(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	CreateUser(&user)
	defer DeleteUser(&user)

	_, err := FindUserByUsername("bob")
	e := err.(*errors.ValidationError)
	msg := "User not found"
	c.Assert(e.Message, Equals, msg)
}
