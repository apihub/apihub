package account_test

import (
	"github.com/backstage/apimanager/account"
	"github.com/backstage/apimanager/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateService(c *C) {
	err := service.Create(owner, team)
	c.Assert(err, IsNil)
	defer service.Delete(owner)
}

func (s *S) TestCreateServiceWithDuplicateAlias(c *C) {
	err := service.Create(owner, team)
	c.Check(err, IsNil)

	err = service.Create(owner, team)
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
	defer service.Delete(owner)
}

func (s *S) TestCreateServiceWithoutRequiredFields(c *C) {
	service = account.Service{}
	err := service.Create(owner, team)
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestUpdateService(c *C) {
	err := service.Create(owner, team)
	c.Assert(err, IsNil)
	c.Assert(service.Endpoint, Equals, "http://example.org/api")

	service.Endpoint = "http://another.org"
	err = service.Update()
	c.Assert(err, IsNil)
	c.Assert(service.Endpoint, Equals, "http://another.org")
	defer service.Delete(owner)
}

func (s *S) TestUpdateServiceWithoutRequiredFields(c *C) {
	err := service.Create(owner, team)
	c.Assert(err, IsNil)
	c.Assert(service.Endpoint, Equals, "http://example.org/api")
	defer service.Delete(owner)

	service.Endpoint = ""
	err = service.Update()
	_, ok := err.(errors.ValidationErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestServiceExists(c *C) {
	service.Create(owner, team)
	c.Assert(service.Exists(), Equals, true)
	defer service.Delete(owner)
}

func (s *S) TestServiceExistsNotFound(c *C) {
	service = account.Service{Subdomain: "not_found"}
	c.Assert(service.Exists(), Equals, false)
}

func (s *S) TestDeleteService(c *C) {
	service.Create(owner, team)
	c.Assert(service.Exists(), Equals, true)
	service.Delete(owner)
	c.Assert(service.Exists(), Equals, false)
}

func (s *S) TestDeleteServiceNotOwner(c *C) {
	service.Create(alice, team)
	c.Assert(service.Exists(), Equals, true)
	defer service.Delete(alice)

	err := service.Delete(owner)
	_, ok := err.(errors.ForbiddenErrorNEW)
	c.Assert(ok, Equals, true)
}

func (s *S) TestFindServiceBySubdomain(c *C) {
	err := service.Create(owner, team)

	t, err := account.FindServiceBySubdomain(service.Subdomain)
	c.Check(t, Not(IsNil))
	c.Check(err, IsNil)
	defer service.Delete(owner)
}

func (s *S) TestFindServiceBySubdomainNotFound(c *C) {
	t, err := account.FindServiceBySubdomain("not-found")
	c.Check(t, IsNil)
	_, ok := err.(errors.NotFoundErrorNEW)
	c.Assert(ok, Equals, true)
}
