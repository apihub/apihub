package account_new_test

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var alice = account_new.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}

func (s *S) TestCreateUser(c *C) {
	defer alice.Delete()
	err := alice.Save()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := account_new.User{}
	err := user.Save()
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestUserToString(c *C) {
	user := account_new.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	str := user.ToString()
	c.Assert(str, Equals, `{"name":"Alice","email":"alice@example.org","username":"alice"}`)
}

func (s *S) TestChangePassword(c *C) {
	alice.Save()
	p1 := alice.Password
	alice.Password = "654321"
	err := alice.ChangePassword()
	defer alice.Delete()
	c.Assert(err, IsNil)
	p2 := alice.Password
	c.Assert(p1, Not(Equals), p2)
}

func (s *S) TestExists(c *C) {
	defer alice.Delete()
	alice.Save()

	valid := alice.Exists()
	c.Assert(valid, Equals, true)
}

func (s *S) TestExistsWhenUserDoesNotExistInTheDB(c *C) {
	valid := alice.Exists()
	c.Assert(valid, Equals, false)
}
