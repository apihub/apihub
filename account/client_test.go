package account

import (
  . "gopkg.in/check.v1"
  "github.com/backstage/backstage/errors"
)

func (s *S) TestCreateClient(c *C) {
  owner := &User{Email: "owner@example.org"}
  team := &Team{Name: "Team", Alias: "team"}
  client := Client{
    Name:  "Backstage App.",
  }
  err := client.Save(owner, team)
  defer client.Delete()
  c.Assert(err, IsNil)
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
  client := &Client{ }
  err := client.Delete()

  e, ok := err.(*errors.ValidationError)
  c.Assert(ok, Equals, true)
  c.Assert(e.Payload, Equals, "Client not found.")
}

func (s *S) TestFindClientByNameAndTeam(c *C) {
  owner := &User{Email: "owner@example.org"}
  team := &Team{Name: "Team", Alias: "team"}
  client := &Client{
    Name: "backstage",
    Team:  team.Alias,
  }

  defer client.Delete()
  client.Save(owner, team)
  se, _ := FindClientByNameAndTeam(client.Name, team.Alias)
  c.Assert(se.Name, Equals, client.Name)
}

func (s *S) TestFindClientByNameAndTeamWithInvalidName(c *C) {
  _, err := FindClientByNameAndTeam("Non Existing Client", "Invalid Team")
  c.Assert(err, NotNil)
  e := err.(*errors.ValidationError)
  c.Assert(e.Payload, Equals, "Client not found.")
}

func (s *S) TestDeleteClientByNameAndTeam(c *C) {
  owner := &User{Email: "owner@example.org"}
  team := &Team{Name: "Team", Alias: "team"}
  client := &Client{
    Name: "backstage",
    Team:  team.Alias,
  }
  defer client.Delete()
  err := client.Save(owner, team)
  c.Assert(err, IsNil)
  err = DeleteClientByNameAndTeam(client.Name, team.Alias)
  c.Assert(err, IsNil)
}

func (s *S) TestDeleteClientByNameAndTeamWithInvalidNameAndTeam(c *C) {
  err := DeleteClientByNameAndTeam("Non existing client", "Invalid Team")
  c.Assert(err, NotNil)
  e := err.(*errors.ValidationError)
  c.Assert(e.Payload, Equals, "Client not found.")
}

func (s *S) TestFindClientsByTeam(c *C) {
  owner := &User{Email: "owner@example.org"}
  team := &Team{Name: "Team", Alias: "team"}
  client := &Client{
    Name: "backstage",
    Team:  team.Alias,
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
