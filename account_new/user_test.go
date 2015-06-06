package account_new_test

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var alice = account_new.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}

func (s *S) TestCreateUser(c *C) {
	defer alice.Delete()
	err := alice.Create()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := account_new.User{}
	err := user.Create()
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestCreateUserWithDuplicateEmail(c *C) {
	defer alice.Delete()
	err := alice.Create()
	c.Check(err, IsNil)

	err = alice.Create()
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestUserToString(c *C) {
	user := account_new.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	str := user.ToString()
	c.Assert(str, Equals, `{"name":"Alice","email":"alice@example.org"}`)
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
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestExists(c *C) {
	defer alice.Delete()
	alice.Create()

	valid := alice.Exists()
	c.Assert(valid, Equals, true)
}

func (s *S) TestExistsWhenUserDoesNotExistInTheDB(c *C) {
	valid := alice.Exists()
	c.Assert(valid, Equals, false)
}
