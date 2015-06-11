package account_new_test

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	err := team.Create(owner)
	defer team.Delete(owner)
	c.Assert(err, IsNil)
}

func (s *S) TestCreateTeamWithDuplicateAlias(c *C) {
	err := team.Create(owner)
	defer team.Delete(owner)
	c.Check(err, IsNil)

	err = team.Create(owner)
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestCreateTeamWithoutRequiredFields(c *C) {
	team = account_new.Team{}
	err := team.Create(owner)
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestTeamExists(c *C) {
	team.Create(owner)
	defer team.Delete(owner)
	c.Assert(team.Exists(), Equals, true)
}

func (s *S) TestTeamExistsNotFound(c *C) {
	team = account_new.Team{Name: "not_found"}
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
	defer team.Delete(owner)

	t, err := account_new.FindTeamByAlias(team.Alias)
	c.Check(t, Not(IsNil))
	c.Check(err, IsNil)
}

func (s *S) TestFindTeamByAliasNotFound(c *C) {
	t, err := account_new.FindTeamByAlias("not-found")
	c.Check(t, IsNil)
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}
