package test

import (
	"testing"

	"github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

var user *account.User
var expires int

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type StorableSuite struct {
	Storage account.Storable
}

func (s *StorableSuite) SetUpTest(c *C) {
	expires = 10
	user = &account.User{Name: "Alice", Username: "alice", Email: "alice@example.org", Password: "123456"}
}

func (s *StorableSuite) TestSaveToken(c *C) {
	key := account.TokenKey{Name: "key"}
	c.Check(s.Storage.SaveToken(key, expires, user), IsNil)
}

func (s *StorableSuite) TestGetToken(c *C) {
	key := account.TokenKey{Name: "keys"}
	s.Storage.SaveToken(key, expires, user)
	var u account.User
	err := s.Storage.GetToken(key, &u)
	c.Check(err, IsNil)
	c.Assert(u.Email, Equals, user.Email)
}

func (s *StorableSuite) TestGetTokenWithNonExistingKey(c *C) {
	key := account.TokenKey{Name: "Non-Existing-Key"}
	var u account.User
	err := s.Storage.GetToken(key, &u)
	c.Assert(u, Equals, account.User{})
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteToken(c *C) {
	key := account.TokenKey{Name: "keys"}
	s.Storage.SaveToken(key, expires, user)
	err := s.Storage.DeleteToken(key)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTokenWithNonExistingKey(c *C) {
	key := account.TokenKey{Name: "Non-Existing-Key"}
	err := s.Storage.DeleteToken(key)
	c.Check(err, IsNil)
}
