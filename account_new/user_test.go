package account_new_test

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	err := alice.Create()
	defer alice.Delete()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := account_new.User{}
	err := user.Create()
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestCreateUserWithDuplicateEmail(c *C) {
	err := alice.Create()
	defer alice.Delete()
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

func (s *S) TestDeleteUser(c *C) {
	alice.Create()
	c.Assert(alice.Exists(), Equals, true)
	alice.Delete()
	c.Assert(alice.Exists(), Equals, false)
}
