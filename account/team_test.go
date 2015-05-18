package account

import (
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName(team.Name)
	c.Assert(err, IsNil)
}

func (s *S) TestSaveExistingTeam(c *C) {
	t := &Team{Name: "My Team", Alias: "myteam"}
	err := t.Save(owner)
	c.Assert(err, IsNil)
	defer DeleteTeamByName(t.Name)

	t, _ = FindTeamByAlias(t.Alias, owner)
	t.Name = "New Name"
	t.Save(owner)
	defer DeleteTeamByName(t.Name)
	c.Assert(t.Name, Equals, "New Name")
}

func (s *S) TestCreateTeamWithoutRequiredFields(c *C) {
	team := Team{}
	err := team.Save(owner)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Name cannot be empty.")
}

func (s *S) TestCreateTeamWhenAliasAlreadyExists(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName(team.Name)
	c.Assert(err, IsNil)

	team = &Team{Name: "Another Team Name", Alias: team.Alias}
	err = team.Save(owner)
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Someone already has that team alias. Could you try another?")
}

func (s *S) TestDeleteTeam(c *C) {
	teamName := team.Name
	team.Save(owner)
	g, _ := FindTeamByName(team.Name)
	c.Assert(len(g.Users), Equals, 1)
	DeleteTeamByAlias(team.Alias, owner)
	_, err := FindTeamByName(teamName)
	c.Assert(err, NotNil)
}

func (s *S) TestAddUsersWithInvalidUser(c *C) {
	err := team.Save(owner)
	defer DeleteTeamByName(team.Name)
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
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "It is not possible to remove the owner from the team.")
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
	c.Assert(e.Payload, Equals, "Team not found.")
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
	c.Assert(e.Payload, Equals, "Team not found.")
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
	c.Assert(e.Payload, Equals, "Team not found.")
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

func (s *S) TestFindTeamByAliasWithServicesAndClients(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	client.Save(owner, team)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()
	defer client.Delete()

	t, _ := FindTeamByAlias(team.Alias, owner)
	c.Assert(t.Name, Equals, "Team")
	c.Assert(t.Services[0], DeepEquals, service)
	c.Assert(t.Clients[0], DeepEquals, client)
}

func (s *S) TestFindTeamByNameWithServicesAndClients(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	client.Save(owner, team)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()
	defer client.Delete()

	t, _ := FindTeamByName(team.Name)
	c.Assert(t.Name, Equals, "Team")
	c.Assert(t.Services[0], DeepEquals, service)
	c.Assert(t.Clients[0], DeepEquals, client)
}

func (s *S) TestFindTeamByIdWithServicesAndClients(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	client.Save(owner, team)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()
	defer client.Delete()

	team, _ = FindTeamByName(team.Name)
	t, _ := FindTeamById(team.Id.Hex())
	c.Assert(t.Services[0], DeepEquals, service)
	c.Assert(t.Clients[0], DeepEquals, client)
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
	c.Assert(err, NotNil)
}

func (s *S) TestTeamToString(c *C) {
	team := &Team{Id: "123", Name: "Backstage", Owner: "alice@example.org", Users: []string{"alice@example.org"}}
	str := team.ToString()
	c.Assert(str, Equals, `{"name":"Backstage","alias":"","users":["alice@example.org"],"owner":"alice@example.org"}`)
}
