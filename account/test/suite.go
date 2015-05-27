package test

import (
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var user *account.User

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type StorableSuite struct {
	Storage account.Storable
}

func (s *StorageSuite) SetUpTest(c *C) {
	user = &account.User{Name: "Alice", Username: "alice", Email: "alice@example.org", Password: "123456"}
}

func (s *StorageSuite) TestSaveToken(c *C) {
	key := account.TokenKey{Name: "key"}
	c.Check(s.Storage.SaveToken(key, user), IsNil)
}

func (s *StorageSuite) TestGetToken(c *C) {
	key := account.TokenKey{Name: "keys"}
	s.Storage.SaveToken(key, user)
	u, err := s.Storage.GetToken(key)
	c.Check(err, IsNil)
	c.Assert(u, Equals, user)
}

func (s *StorageSuite) TestGetTokenWithNonExistingKey(c *C) {
	key := account.TokenKey{Name: "Non-Existing-Key"}
	u, err := s.Storage.GetToken(key)
	c.Check(u, IsNil)
	e := err.(*errors.NotFoundError)
	c.Check(e, Not(IsNil))
}

func (s *StorageSuite) TestDeleteToken(c *C) {
	key := account.TokenKey{Name: "keys"}
	s.Storage.SaveToken(key, user)
	err := s.Storage.DeleteToken(key)
	c.Check(err, IsNil)
}

func (s *StorageSuite) TestDeleteTokenWithNonExistingKey(c *C) {
	key := account.TokenKey{Name: "Non-Existing-Key"}
	err := s.Storage.DeleteToken(key)
	e := err.(*errors.NotFoundError)
	c.Check(e, Not(IsNil))
}
