package account

import (
	. "gopkg.in/check.v1"

	"github.com/albertoleal/backstage/errors"
)

func (s *S) TestCreateUser(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	err := user.Save()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithSameUsername(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	user2 := User{Name: "Bob", Email: "bob@bar.com", Username: "alice", Password: "123456"}
	err := user2.Save()
	e := err.(*errors.ValidationError)
	msg := "Someone already has that username. Could you try another?."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := User{}
	err := user.Save()
	e := err.(*errors.ValidationError)
	msg := "Name/Email/Username/Password cannot be empty."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCreateUserShouldMaskThePassword(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	foundUser, _ := FindUserByUsername("alice")
	c.Assert(foundUser.Password, Not(Equals), "123456")
}

func (s *S) TestValid(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	valid := user.Valid()
	c.Assert(valid, Equals, true)
}

func (s *S) TestValidWhenUserDoesNotExistInTheDB(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	valid := user.Valid()
	c.Assert(valid, Equals, false)
}

func (s *S) FindUserByUsername(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	foundUser, err := FindUserByUsername("alice")
	c.Assert(err, IsNil)
	c.Assert(foundUser, NotNil)
}

func (s *S) TestFindUserWithInvalidUsername(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	_, err := FindUserByUsername("bob")
	e := err.(*errors.ValidationError)
	msg := "User not found"
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestGetTeams(c *C) {
	user := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()
	group := &Group{Name: "Group"}
	group.Save(&user)
	defer group.Delete()
	g, err := user.GetTeams()
	c.Assert(err, IsNil)
	c.Assert(g[0].Name, Equals, "Group")
}
