package app

import (
	"github.com/albertoleal/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateGroup(c *C) {
	err := CreateGroup("Group", []User{})
	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)
}

func (s *S) TestCreateGroupWhenNameAlreadyExists(c *C) {
	err := CreateGroup("Group", []User{})
	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)

	err = CreateGroup("Group", []User{})
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Someone already has that group name. Could you try another?"
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestAddUsers(c *C) {
	err := CreateGroup("Group", []User{})
	defer DeleteGroupByName("Group")
	group, _ := FindGroupByName("Group")

	alice := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	err = group.AddUsers([]User{alice, bob})
	c.Assert(err, IsNil)

	group, _ = FindGroupByName("Group")
	c.Assert(len(group.Users), Equals, 2)
}

func (s *S) TestAddUsersWithSameUsername(c *C) {
	err := CreateGroup("Group", []User{})
	defer DeleteGroupByName("Group")
	group, _ := FindGroupByName("Group")
	c.Assert(len(group.Users), Equals, 0)

	alice := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "alice", Password: "123456"}
	err = group.AddUsers([]User{alice, bob})
	c.Assert(err, IsNil)

	group, _ = FindGroupByName("Group")
	c.Assert(len(group.Users), Equals, 1)
}

func (s *S) TestRemoveUsers(c *C) {
	err := CreateGroup("Group", []User{})
	defer DeleteGroupByName("Group")
	group, _ := FindGroupByName("Group")

	alice := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	group.AddUsers([]User{alice, bob})

	err = group.RemoveUsers([]User{bob})
	c.Assert(err, IsNil)

	group, _ = FindGroupByName("Group")
	c.Assert(len(group.Users), Equals, 1)
	c.Assert(group.Users[0], Equals, "alice")
}

func (s *S) TestRemoveUsersWithNonExistingUser(c *C) {
	err := CreateGroup("Group", []User{})
	defer DeleteGroupByName("Group")
	group, _ := FindGroupByName("Group")

	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	err = group.RemoveUsers([]User{bob})
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteGroupByName(c *C) {
	err := CreateGroup("Group", []User{})
	c.Assert(err, IsNil)
	err = DeleteGroupByName("Group")
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteGroupByNameWithInvalidName(c *C) {
	err := DeleteGroupByName("Non Existing Group")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Group not found."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestFindGroupByName(c *C) {
	alice := User{Name: "Alice", Email: "alice@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	err := CreateGroup("Group", []User{alice, bob})
	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)

	group, _ := FindGroupByName("Group")
	c.Assert(len(group.Users), Equals, 2)
	c.Assert(group.Name, Equals, "Group")
}

func (s *S) TestFindGroupByNameWithInvalidName(c *C) {
	_, err := FindGroupByName("Non Existing Group")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Group not found."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestGetGroupUsers(c *C) {
	alice := User{Name: "Alice", Email: "alice@bar.com", Username: "alice", Password: "123456"}
	CreateUser(&alice)
	defer DeleteUser(&alice)
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	CreateUser(&bob)
	defer DeleteUser(&bob)
	CreateGroup("Group", []User{alice, bob})
	defer DeleteGroupByName("Group")

	group, _ := FindGroupByName("Group")
	users, _ := group.GetGroupUsers()
	c.Assert(users[0], DeepEquals, &alice)
	c.Assert(users[1], DeepEquals, &bob)
}
