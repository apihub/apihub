package account_new_test

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	err := team.Create(owner)
	defer team.Delete()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateTeamWithDuplicateAlias(c *C) {
	err := team.Create(owner)
	defer team.Delete()
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
	defer team.Delete()
	c.Assert(team.Exists(), Equals, true)
}

func (s *S) TestTeamExistsNotFound(c *C) {
	c.Assert(team.Exists(), Equals, false)
}

func (s *S) TestDeleteTeam(c *C) {
	team.Create(owner)
	c.Assert(team.Exists(), Equals, true)
	team.Delete()
	c.Assert(team.Exists(), Equals, false)
}
