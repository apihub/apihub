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

func (s *S) TestDeleteGroup(c *C) {
	groupName := group.Name
	group.Save(owner)
	g, _ := FindGroupByName(group.Name)
	c.Assert(len(g.Users), Equals, 1)
	group.Delete()
	_, err := FindGroupByName(groupName)
	c.Assert(err, Not(IsNil))
}

func (s *S) TestAddUsersWithInvalidUser(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")

	err = g.AddUsers([]string{"alice", "bob"})
	c.Assert(err, IsNil)

	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)
}

func (s *S) TestAddUsersWithValidUser(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")

	bob := &User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	err = g.AddUsers([]string{"alice", "bob"})
	c.Assert(err, IsNil)

	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 2)
}

func (s *S) TestAddUsersWithSameUsername(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)

	err = g.AddUsers([]string{"alice", "alice"})
	c.Assert(err, IsNil)

	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)
}

func (s *S) TestRemoveUsers(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")
	err = g.AddUsers([]string{"alice", "bob"})
	err = g.RemoveUsers([]string{"bob"})
	c.Assert(err, IsNil)
	g, _ = FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 1)
	c.Assert(g.Users[0], Equals, "alice")
}

func (s *S) TestRemoveUsersWithNonExistingUser(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	g, _ := FindGroupByName("Group")
	err = g.RemoveUsers([]string{"bob"})
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveUsersWhenTheUserIsOwner(c *C) {
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	mary := &User{Name: "Mary", Email: "mary@bar.com", Username: "mary", Password: "123456"}
	mary.Save()
	defer mary.Delete()
	group.AddUsers([]string{"mary", "bob"})

	err = group.RemoveUsers([]string{owner.Username, "bob"})
	c.Assert(err, Not(IsNil))
	e := err.(*errors.ValidationError)
	c.Assert(e.Message, Equals, "It is not possible to remove the owner from the team.")

	g, _ := FindGroupByName("Group")
	c.Assert(len(g.Users), Equals, 2)
	c.Assert(g.Users[0], Equals, owner.Username)
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
	owner := &User{Name: "Alice", Email: "alice@bar.com", Username: "alice", Password: "123456"}
	err := group.Save(owner)

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

func (s *S) TestFindGroupById(c *C) {
	owner := &User{Name: "Alice", Email: "alice@bar.com", Username: "alice", Password: "123456"}
	err := group.Save(owner)
	defer DeleteGroupByName("Group")
	c.Assert(err, IsNil)
	group, _ := FindGroupByName("Group")
	g, _ := FindGroupById(group.Id.Hex())
	c.Assert(g.Name, Equals, "Group")
}

func (s *S) TestFindGroupByIdWithInvalidId(c *C) {
	_, err := FindGroupById("123")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Group not found."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestGetGroupUsers(c *C) {
	alice := &User{Name: "Alice", Email: "alice@bar.com", Username: "alice", Password: "123456"}
	defer alice.Delete()
	alice.Save()
	bob := &User{Name: "Bob", Email: "bob@bar.com", Username: "bob", Password: "123456"}
	defer bob.Delete()
	bob.Save()

	group.Save(alice)
	group.AddUsers([]string{"bob"})
	defer DeleteGroupByName("Group")

	g, _ := FindGroupByName("Group")
	users, _ := g.GetGroupUsers()
	c.Assert(users[0].Username, Equals, alice.Username)
	c.Assert(users[1].Username, Equals, bob.Username)
}
