package test

import (
	"testing"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var user account_new.User
var expires int

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type StorableSuite struct {
	Storage account_new.Storable
}

func (s *StorableSuite) SetUpTest(c *C) {
	expires = 10
	user = account_new.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
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
