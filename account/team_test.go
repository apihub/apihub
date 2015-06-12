package account_test

import (
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	err := team.Create(owner)
	c.Assert(err, IsNil)
	defer team.Delete(owner)
}

func (s *S) TestCreateTeamWithDuplicateAlias(c *C) {
	err := team.Create(owner)
	c.Check(err, IsNil)

	err = team.Create(owner)
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
	defer team.Delete(owner)
}

func (s *S) TestCreateTeamWithoutRequiredFields(c *C) {
	team = account.Team{}
	err := team.Create(owner)
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestUpdateTeam(c *C) {
	err := team.Create(owner)
	c.Assert(err, IsNil)
	c.Assert(team.Name, Equals, "Backstage Team")

	team.Name = "New name"
	err = team.Update()
	c.Assert(err, IsNil)
	c.Assert(team.Name, Equals, "New name")
	defer team.Delete(owner)
}

func (s *S) TestUpdateTeamWithoutRequiredFields(c *C) {
	err := team.Create(owner)
	c.Assert(err, IsNil)
	c.Assert(team.Name, Equals, "Backstage Team")

	team.Name = ""
	err = team.Update()
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
	defer team.Delete(owner)
}

func (s *S) TestTeamExists(c *C) {
	team.Create(owner)
	c.Assert(team.Exists(), Equals, true)
	defer team.Delete(owner)
}

func (s *S) TestTeamExistsNotFound(c *C) {
	team = account.Team{Name: "not_found"}
	c.Assert(team.Exists(), Equals, false)
}

func (s *S) TestDeleteTeam(c *C) {
	team.Create(owner)
	c.Assert(team.Exists(), Equals, true)
	team.Delete(owner)
	c.Assert(team.Exists(), Equals, false)
}

func (s *S) TestDeleteTeamNotOwner(c *C) {
	team.Create(alice)
	c.Assert(team.Exists(), Equals, true)
	err := team.Delete(owner)
	_, ok := err.(errors.ForbiddenErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestFindTeamByAlias(c *C) {
	err := team.Create(owner)

	t, err := account.FindTeamByAlias(team.Alias)
	c.Check(t, Not(IsNil))
	c.Check(err, IsNil)
	defer team.Delete(owner)
}

func (s *S) TestFindTeamByAliasNotFound(c *C) {
	t, err := account.FindTeamByAlias("not-found")
	c.Check(t, IsNil)
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestContainsUser(c *C) {
	team.Users = append(team.Users, alice.Email)
	pos, err := team.ContainsUser(&alice)
	c.Check(err, IsNil)
	c.Assert(pos >= 0, Equals, true)
}

func (s *S) TestContainsUserNotFound(c *C) {
	pos, err := team.ContainsUser(&alice)
	_, ok := err.(errors.ForbiddenErrorNEW)
	c.Assert(ok, Equals, true)
	c.Assert(pos, Equals, -1)
}

func (s *S) TestAddUsers(c *C) {
	team.Create(owner)
	alice.Create()
	err := team.AddUsers([]string{alice.Email})
	c.Check(err, IsNil)
	defer alice.Delete()
	defer team.Delete(owner)
}

func (s *S) TestAddUsersWithInvalidUser(c *C) {
	team.Create(owner)
	err := team.AddUsers([]string{"maria@example.org"})
	c.Assert(err, IsNil)

	t, err := account.FindTeamByAlias(team.Alias)
	c.Assert(len(t.Users), Equals, 1)
	defer team.Delete(owner)
}

func (s *S) TestAddUsersWithSameUsername(c *C) {
	team.Create(owner)
	err := team.AddUsers([]string{owner.Email})
	c.Assert(err, IsNil)

	t, err := account.FindTeamByAlias(team.Alias)
	c.Assert(len(t.Users), Equals, 1)
	defer team.Delete(owner)
}

func (s *S) TestRemoveUsers(c *C) {
	team.Create(owner)
	alice.Create()

	err := team.AddUsers([]string{alice.Email})
	c.Check(err, IsNil)
	t, err := account.FindTeamByAlias(team.Alias)
	c.Assert(len(t.Users), Equals, 2)

	err = team.RemoveUsers([]string{alice.Email})
	c.Check(err, IsNil)
	t, err = account.FindTeamByAlias(team.Alias)
	c.Assert(len(t.Users), Equals, 1)

	defer alice.Delete()
	defer team.Delete(owner)
}

func (s *S) TestRemoveUsersWithInvalidUser(c *C) {
	team.Create(owner)
	alice.Create()

	err := team.RemoveUsers([]string{"invalid@example.org"})
	c.Check(err, IsNil)
	t, err := account.FindTeamByAlias(team.Alias)
	c.Assert(len(t.Users), Equals, 1)

	defer alice.Delete()
	defer team.Delete(owner)
}

func (s *S) TestRemoveUsersWhenTheUserIsOwner(c *C) {
	team.Create(owner)
	owner.Create()

	err := team.RemoveUsers([]string{owner.Email})
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)

	t, err := account.FindTeamByAlias(team.Alias)
	c.Assert(len(t.Users), Equals, 1)

	defer team.Delete(owner)
	defer owner.Delete()
}
