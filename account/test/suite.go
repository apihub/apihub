package test

import (
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var user account.User
var expires int

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type StorableSuite struct {
	Storage account.Storable
}

func (s *StorableSuite) SetUpTest(c *C) {
	expires = 10
	user = account.User{Name: "Alice", Username: "alice", Email: "alice@example.org", Password: "123456"}
}

func (s *StorableSuite) TestCreateUser(c *C) {
	defer s.Storage.DeleteUser(user)
	err := s.Storage.CreateUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestCreateUserWithDupEmail(c *C) {
	user.Username = "alice1"
	defer s.Storage.DeleteUser(user)
	err := s.Storage.CreateUser(user)
	c.Check(err, IsNil)

	user.Username = "alice2"
	err = s.Storage.CreateUser(user)
	c.Assert(err, Equals, errors.ErrUserDuplicateEntry)
}

func (s *StorableSuite) TestCreateUserWithDupUsername(c *C) {
	user.Email = "alice@example.org"
	defer s.Storage.DeleteUser(user)
	err := s.Storage.CreateUser(user)
	c.Check(err, IsNil)

	user.Email = "alice2@example.org"
	err = s.Storage.CreateUser(user)
	c.Assert(err, Equals, errors.ErrUserDuplicateEntry)
}

func (s *StorableSuite) TestUpdateUser(c *C) {
	s.Storage.CreateUser(user)
	user.Name = "Bob"
	defer s.Storage.DeleteUser(user)
	err := s.Storage.UpdateUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestUpdateUserNotFound(c *C) {
	err := s.Storage.UpdateUser(user)
	c.Assert(err, Equals, errors.ErrUserNotFound)
}

func (s *StorableSuite) TestDeleteUser(c *C) {
	s.Storage.CreateUser(user)
	err := s.Storage.DeleteUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteUserNotFound(c *C) {
	err := s.Storage.DeleteUser(user)
	c.Assert(err, Equals, errors.ErrUserNotFound)
}

func (s *StorableSuite) TestFindUserByEmail(c *C) {
	defer s.Storage.DeleteUser(user)
	s.Storage.CreateUser(user)
	u, err := s.Storage.FindUserByEmail(user.Email)
	c.Assert(u, Equals, user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindUserByEmailNotFound(c *C) {
	_, err := s.Storage.FindUserByEmail("not-found")
	c.Assert(err, Equals, errors.ErrUserNotFound)
}
