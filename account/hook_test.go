package account_test

import (
	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestSaveHook(c *C) {
	err := hook.Save(team)
	c.Assert(err, IsNil)
	defer hook.Delete()
}

func (s *S) TestSaveHookWithoutRequiredFields(c *C) {
	hook = account.Hook{}
	err := hook.Save(team)
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestDeleteHook(c *C) {
	hook.Save(team)
	c.Assert(hook.Exists(), Equals, true)
	hook.Delete()
	c.Assert(hook.Exists(), Equals, false)
}

func (s *S) TestFindHookByName(c *C) {
	hook.Save(team)

	t, err := account.FindHookByName(hook.Name)
	c.Check(t, Not(IsNil))
	c.Check(err, IsNil)
	defer hook.Delete()
}

func (s *S) TestFindHookByNameNotFound(c *C) {
	t, err := account.FindHookByName("not-found")
	c.Check(t, IsNil)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}
