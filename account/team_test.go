package account

import (
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var (
	team  *Team
	owner *User
)

func (s *S) SetUpTest(c *C) {
	team = &Team{Name: "Team", Alias: "Alias"}
	owner = &User{Name: "Owner", Username: "owner", Email: "owner@example.org", Password: "123456"}
}

func (s *S) TestCreateTeam(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	c.Assert(err, IsNil)
}

func (s *S) TestCreateTeamWithoutRequiredFields(c *C) {
	team := Team{}
	err := team.Save(owner)
	e := err.(*errors.ValidationError)
	msg := "Name cannot be empty."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCreateTeamWhenAliasAlreadyExists(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	c.Assert(err, IsNil)

	team = &Team{Name: "Another Team Name", Alias: team.Alias}
	err = team.Save(owner)
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Someone already has that team alias. Could you try another?"
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestDeleteTeam(c *C) {
	teamName := team.Name
	team.Save(owner)
	g, _ := FindTeamByName(team.Name)
	c.Assert(len(g.Users), Equals, 1)
	DeleteTeamByAlias(team.Alias, owner)
	_, err := FindTeamByName(teamName)
	c.Assert(err, Not(IsNil))
}

func (s *S) TestAddUsersWithInvalidUser(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	g, _ := FindTeamByName("Team")

	err = g.AddUsers([]string{"owner@example.org", "bob@example.org"})
	c.Assert(err, IsNil)

	g, _ = FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 1)
}

func (s *S) TestAddUsersWithValidUser(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	g, _ := FindTeamByName("Team")

	bob := &User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	err = g.AddUsers([]string{bob.Email})
	c.Assert(err, IsNil)

	g, _ = FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 2)
}

func (s *S) TestAddUsersWithSameUsername(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	g, _ := FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 1)
	err = g.AddUsers([]string{"alice@example.org", "alice@example.org"})
	c.Assert(err, IsNil)
	g, _ = FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 1)
}

func (s *S) TestRemoveUsers(c *C) {
	err := team.Save(owner)
	bob := &User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	defer DeleteTeamByName("Team")
	g, _ := FindTeamByName("Team")
	err = g.AddUsers([]string{bob.Email})
	g, _ = FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 2)
	err = g.RemoveUsers([]string{bob.Email})
	c.Assert(err, IsNil)
	g, _ = FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 1)
	c.Assert(g.Users[0], Equals, "owner@example.org")
}

func (s *S) TestRemoveUsersWithNonExistingUser(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	g, _ := FindTeamByName("Team")
	err = g.RemoveUsers([]string{"bob@example.org"})
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveUsersWhenTheUserIsOwner(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	mary := &User{Name: "Mary", Email: "mary@example.org", Username: "mary", Password: "123456"}
	mary.Save()
	defer mary.Delete()
	team.AddUsers([]string{"mary@example.org", "bob@example.org"})

	err = team.RemoveUsers([]string{owner.Email, "bob@example.org"})
	c.Assert(err, Not(IsNil))
	e := err.(*errors.ValidationError)
	c.Assert(e.Message, Equals, "It is not possible to remove the owner from the team.")

	g, _ := FindTeamByName("Team")
	c.Assert(len(g.Users), Equals, 2)
	c.Assert(g.Users[0], Equals, owner.Email)
}

func (s *S) TestDeleteTeamByName(c *C) {
	err := team.Save(owner)
	c.Assert(err, IsNil)
	err = DeleteTeamByName("Team")
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteTeamByNameWithInvalidName(c *C) {
	err := DeleteTeamByName("Non Existing Team")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Team not found."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestFindTeamByName(c *C) {
	err := team.Save(owner)

	defer DeleteTeamByName("Team")
	c.Assert(err, IsNil)

	g, _ := FindTeamByName("Team")
	c.Assert(g.Name, Equals, "Team")
}

func (s *S) TestFindTeamByNameWithInvalidName(c *C) {
	_, err := FindTeamByName("Non Existing Team")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Team not found."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestFindTeamById(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	c.Assert(err, IsNil)
	team, _ := FindTeamByName("Team")
	g, _ := FindTeamById(team.Id.Hex())
	c.Assert(g.Name, Equals, "Team")
}

func (s *S) TestFindTeamByIdWithInvalidId(c *C) {
	_, err := FindTeamById("123")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	message := "Team not found."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestFindTeamByAlias(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName("Team")
	c.Assert(err, IsNil)

	g, _ := FindTeamByAlias("alias", owner)
	c.Assert(g.Name, Equals, "Team")
}

func (s *S) TestFindTeamByAliasWithInvalidName(c *C) {
	_, err := FindTeamByAlias("Non Existing Team", owner)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Team not found.")
}

func (s *S) TestGetTeamUsers(c *C) {
	alice := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer alice.Delete()
	alice.Save()
	bob := &User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	defer bob.Delete()
	bob.Save()

	team.Save(alice)
	team.AddUsers([]string{bob.Email})
	defer DeleteTeamByName("Team")

	g, _ := FindTeamByName("Team")
	users, _ := g.GetTeamUsers()
	c.Assert(users[0].Email, Equals, alice.Email)
	c.Assert(users[1].Email, Equals, bob.Email)
}

func (s *S) TestContainsUser(c *C) {
	bob := &User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	defer bob.Delete()
	bob.Save()

	team.Save(owner)
	defer DeleteTeamByName("Team")
	g, _ := FindTeamByName("Team")
	_, err := g.ContainsUser(owner)
	c.Assert(err, IsNil)
	_, err = g.ContainsUser(bob)
	c.Assert(err, Not(IsNil))
}

func (s *S) TestTeamToString(c *C) {
	team := &Team{Id: "123", Name: "Backstage", Owner: "alice@example.org", Users: []string{"alice@example.org"}}
	str := team.ToString()
	c.Assert(str, Equals, `{"name":"Backstage","alias":"","users":["alice@example.org"],"owner":"alice@example.org"}`)
}
