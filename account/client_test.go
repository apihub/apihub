package account

import (
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateClient(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := Client{
		Name: "Backstage App.",
	}
	err := client.Save(owner, team)
	defer client.Delete()
	c.Assert(err, IsNil)
}

func (s *S) TestSaveExistingClient(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := Client{
		Name: "Backstage App.",
	}
	err := client.Save(owner, team)
	defer client.Delete()
	c.Assert(err, IsNil)

	client.Name = "New name"
	err = client.Save(owner, team)
	defer DeleteClientByIdAndTeam(client.Id, team.Alias)

	c.Assert(client.Name, Equals, "New name")
	c.Check(err, IsNil)
}

func (s *S) TestCannotCreateClientWithoutRequiredFields(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{}
	err := client.Save(owner, team)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Name cannot be empty.")
}

func (s *S) TestDeleteClient(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{Name: "backstage"}

	err := client.Save(owner, team)
	c.Assert(err, IsNil)
	err = client.Delete()
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteClientWhenClientDoesNotExist(c *C) {
	client := &Client{}
	err := client.Delete()

	e, ok := err.(*errors.NotFoundError)
	c.Assert(ok, Equals, true)
	c.Assert(e.Payload, Equals, "Client not found.")
}

func (s *S) TestFindClientById(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{
		Id:   "backstage",
		Name: "backstage",
		Team: team.Alias,
	}

	defer client.Delete()
	client.Save(owner, team)
	se, err := FindClientById(client.Id)
	panic(err)
	c.Assert(se.Name, Equals, client.Name)
	c.Assert(err, IsNil)
}

func (s *S) TestFindClientByIdWithInvalidId(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{
		Id:   "backstage",
		Name: "backstage",
		Team: team.Alias,
	}

	defer client.Delete()
	client.Save(owner, team)
	_, err := FindClientById("invalid-id")
	c.Assert(err, NotNil)
}

func (s *S) TestFindClientByIdAndTeam(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{
		Id:   "backstage",
		Name: "backstage",
		Team: team.Alias,
	}

	defer client.Delete()
	client.Save(owner, team)
	se, err := FindClientByIdAndTeam(client.Id, team.Alias)
	c.Assert(se.Name, Equals, client.Name)
	c.Assert(err, IsNil)
}

func (s *S) TestFindClientByIdAndTeamWithInvalidName(c *C) {
	_, err := FindClientByIdAndTeam("Non Existing Client", "Invalid Team")
	c.Assert(err, NotNil)
	e := err.(*errors.NotFoundError)
	c.Assert(e.Payload, Equals, "Client not found.")
}

func (s *S) TestDeleteClientByIdAndTeam(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{
		Name: "backstage",
		Team: team.Alias,
	}
	defer client.Delete()
	err := client.Save(owner, team)
	c.Assert(err, IsNil)
	err = DeleteClientByIdAndTeam(client.Id, team.Alias)
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteClientByIdAndTeamWithInvalidNameAndTeam(c *C) {
	err := DeleteClientByIdAndTeam("Non existing client", "Invalid Team")
	c.Assert(err, NotNil)
	e := err.(*errors.NotFoundError)
	c.Assert(e.Payload, Equals, "Client not found.")
}

func (s *S) TestDeleteClientByTeam(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{
		Name: "backstage",
		Team: team.Alias,
	}
	defer client.Delete()
	err := client.Save(owner, team)
	c.Assert(err, IsNil)
	err = DeleteClientByTeam(team.Alias)
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteClientByTeamWithInvalidTeam(c *C) {
	err := DeleteClientByTeam("Invalid Team")
	c.Assert(err, IsNil)
}

func (s *S) TestFindClientsByTeam(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	client := &Client{
		Name: "backstage",
		Team: team.Alias,
	}

	defer client.Delete()
	client.Save(owner, team)
	se, _ := FindClientsByTeam(team.Alias)
	c.Assert(len(se), Equals, 1)
	c.Assert(se[0].Name, Equals, "backstage")
}

func (s *S) TestFindClientsByTeamWithoutElements(c *C) {
	se, _ := FindClientsByTeam("non-existing-team")
	c.Assert(len(se), Equals, 0)
}
