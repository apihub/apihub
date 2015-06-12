package test

import (
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var user account.User
var team account.Team
var token account.TokenInfo

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type StorableSuite struct {
	Storage account.Storable
}

func (s *StorableSuite) SetUpTest(c *C) {
	user = account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	team = account.Team{Name: "Backstage Team", Alias: "backstage", Users: []string{user.Email}, Owner: user.Email}
	token = account.TokenInfo{Token: "secret-token", Expires: 10, Type: "Token", User: &user}
}

func (s *StorableSuite) TestUpsertUser(c *C) {
	defer s.Storage.DeleteUser(user)
	err := s.Storage.UpsertUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestUpdateUser(c *C) {
	s.Storage.UpsertUser(user)
	user.Name = "Bob"
	defer s.Storage.DeleteUser(user)
	err := s.Storage.UpsertUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestTeams(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)
	teams, err := s.Storage.UserTeams(team.Owner)
	c.Check(err, IsNil)
	c.Assert(teams, DeepEquals, []account.Team{team})
}

func (s *StorableSuite) TestTeamsNotFound(c *C) {
	teams, err := s.Storage.UserTeams("not-found")
	c.Check(err, IsNil)
	c.Assert(teams, DeepEquals, []account.Team{})
}

func (s *StorableSuite) TestDeleteUser(c *C) {
	s.Storage.UpsertUser(user)
	err := s.Storage.DeleteUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteUserNotFound(c *C) {
	err := s.Storage.DeleteUser(user)
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestFindUserByEmail(c *C) {
	defer s.Storage.DeleteUser(user)
	s.Storage.UpsertUser(user)
	u, err := s.Storage.FindUserByEmail(user.Email)
	c.Assert(u, Equals, user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindUserByEmailNotFound(c *C) {
	_, err := s.Storage.FindUserByEmail("not-found")
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestUpsertTeam(c *C) {
	defer s.Storage.DeleteTeam(team)
	err := s.Storage.UpsertTeam(team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTeam(c *C) {
	s.Storage.UpsertTeam(team)
	err := s.Storage.DeleteTeam(team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTeamNotFound(c *C) {
	err := s.Storage.DeleteTeam(team)
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestDeleteTeamByAlias(c *C) {
	s.Storage.UpsertTeam(team)
	err := s.Storage.DeleteTeamByAlias(team.Alias)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTeamByAliasNotFound(c *C) {
	err := s.Storage.DeleteTeamByAlias(team.Alias)
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestFindTeamByAlias(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)
	u, err := s.Storage.FindTeamByAlias(team.Alias)
	c.Assert(u, DeepEquals, team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindTeamByAliasNotFound(c *C) {
	_, err := s.Storage.FindTeamByAlias("not-found")
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestCreateToken(c *C) {
	defer s.Storage.DeleteToken(token.Token)
	err := s.Storage.CreateToken(token)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteToken(c *C) {
	err := s.Storage.CreateToken(token)
	c.Check(err, IsNil)
	err = s.Storage.DeleteToken(token.Token)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDecodeToken(c *C) {
	s.Storage.CreateToken(token)
	var u account.User
	s.Storage.DecodeToken(token.Token, &u)
	c.Assert(u, DeepEquals, user)
}