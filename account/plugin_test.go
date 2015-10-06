package account_test

import (
	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestSavePlugin(c *C) {
	err := pluginConfig.Save(service)
	c.Assert(err, IsNil)
	defer pluginConfig.Delete()
}

func (s *S) TestSavePluginWithoutRequiredFields(c *C) {
	pluginConfig := account.Plugin{}
	err := pluginConfig.Save(service)
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestFindPluginByNameAndService(c *C) {
	err := pluginConfig.Save(service)

	t, err := account.FindPluginByNameAndService(pluginConfig.Name, service)
	c.Check(t, Not(IsNil))
	c.Check(err, IsNil)
	defer service.Delete(owner)
}

func (s *S) TestFindPluginByNameAndServiceNotFound(c *C) {
	t, err := account.FindPluginByNameAndService("not-found", service)
	c.Check(t, IsNil)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}
