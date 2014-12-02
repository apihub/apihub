package account

import (
	"github.com/albertoleal/backstage/errors"
	. "gopkg.in/check.v1"
)

var (
	group *Group
	owner *User
)

func (s *S) SetUpTest(c *C) {
	group = &Group{Name: "Group"}
	owner = &User{Name: "Alice", Username: "alice"}
}

func (s *S) TestCreateGroup(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)
}

func (s *S) TestCreateGroupWhenNameAlreadyExists(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)

	group = &Group{Name: "Group"}
	err = group.Save(owner)
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Someone already has that group name. Could you try another?"
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestAddUsers(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")

	alice := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	err = g.AddUsers([]User{alice, bob})
	c.Assert(err, IsNil)

	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 2)
}

func (s *S) TestAddUsersWithSameUsername(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)

	alice := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "alice", Password: "123456"}
	err = g.AddUsers([]User{alice, bob})
	c.Assert(err, IsNil)

	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)
}

func (s *S) TestRemoveUsers(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")

	alice := User{Name: "Alice", Email: "foo@bar.com", Username: "alice", Password: "123456"}
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	g.AddUsers([]User{alice, bob})

	err = g.RemoveUsers([]User{bob})
	c.Assert(err, IsNil)

	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)
	c.Assert(g.Users[0], Equals, "alice")
}

func (s *S) TestRemoveUsersWithNonExistingUser(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")

	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	err = g.RemoveUsers([]User{bob})
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteGroupByName(c *C) {
	err := group.Save(owner)
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
	alice := &User{Name: "Alice", Email: "alice@bar.com", Username: "alice", Password: "123456"}
	err := group.Save(alice)

	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)

	g, _ := FindGroupByName("Group")
	c.Assert(g.Name, Equals, "Group")
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
	defer alice.Delete()
	alice.Save()
	bob := User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	defer bob.Delete()
	bob.Save()

	group.Save(&alice)
	group.AddUsers([]User{bob})
	defer DeleteGroupByName("Group")

	g, _ := FindGroupByName("Group")
	users, _ := g.GetGroupUsers()
	c.Assert(users[0], DeepEquals, &alice)
	c.Assert(users[1], DeepEquals, &bob)
}
