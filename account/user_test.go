package account_test

import (
	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	err := alice.Create()
	defer alice.Delete()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := account.User{}
	err := user.Create()
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestCreateUserWithDuplicateEmail(c *C) {
	err := alice.Create()
	defer alice.Delete()
	c.Check(err, IsNil)

	err = alice.Create()
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestChangePassword(c *C) {
	alice.Create()
	p1 := alice.Password
	alice.Password = "654321"
	err := alice.ChangePassword()
	defer alice.Delete()
	c.Assert(err, IsNil)
	p2 := alice.Password
	c.Assert(p1, Not(Equals), p2)
}

func (s *S) TestChangePasswordNotFound(c *C) {
	err := alice.ChangePassword()
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestUserExists(c *C) {
	alice.Create()
	defer alice.Delete()

	valid := alice.Exists()
	c.Assert(valid, Equals, true)
}

func (s *S) TestUserExistsNotFound(c *C) {
	valid := alice.Exists()
	c.Assert(valid, Equals, false)
}

func (s *S) TestTeams(c *C) {
	err := team.Create(alice)
	defer team.Delete(alice)
	teams, err := alice.Teams()
	c.Assert(err, IsNil)
	c.Assert(teams, DeepEquals, []account.Team{team})
}

func (s *S) TestDeleteUser(c *C) {
	alice.Create()
	c.Assert(alice.Exists(), Equals, true)
	alice.Delete()
	c.Assert(alice.Exists(), Equals, false)
}

func (s *S) TestServices(c *C) {
	err := team.Create(alice)
	defer team.Delete(alice)
	err = service.Create(alice, team)
	defer service.Delete(alice)

	services, err := alice.Services()
	c.Assert(err, IsNil)
	c.Assert(services, DeepEquals, []account.Service{service})
}
