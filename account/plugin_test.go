package account

import (
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreatePluginConfigNewService(c *C) {
	owner.Save()
	team.Save(owner)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	service := Service{
		Endpoint:  "http://example.org/api",
		Subdomain: "_test_middleware_config_new_service",
	}
	service.Save(owner, team)
	defer DeleteServiceBySubdomain("backstage")

	config := &PluginConfig{
		Name:    "cors",
		Service: service.Subdomain,
		Config:  map[string]interface{}{"allowed_origins": []string{"www.example.org"}, "debug": true},
	}
	defer config.Delete(owner)
	err := config.Save(owner)
	_, ok := err.(*errors.ValidationError)
	c.Check(err, IsNil)
	c.Check(ok, Equals, false)
}

func (s *S) TestCannotCreatePluginConfigWithoutRequiredFields(c *C) {
	owner.Save()
	team.Save(owner)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	midd := &PluginConfig{}
	err := midd.Save(owner)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Name cannot be empty.")

	midd = &PluginConfig{Name: "foo"}
	err = midd.Save(owner)
	e = err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Service cannot be empty.")
}

func (s *S) TestDeleteMiddConfigANonExistingMidd(c *C) {
	owner.Save()
	team.Save(owner)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	config := &PluginConfig{
		Name:    "foo",
		Service: "backstage",
	}
	err := config.Delete(owner)
	c.Check(err, NotNil)
}

func (s *S) TestDeletePluginsByService(c *C) {
	owner.Save()
	team.Save(owner)
	defer DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	service := Service{
		Endpoint:  "http://example.org/api",
		Subdomain: "_test_middleware_config_new_service",
	}
	service.Save(owner, team)
	defer DeleteServiceBySubdomain("backstage")

	config := &PluginConfig{
		Name:    "cors",
		Service: service.Subdomain,
		Config:  map[string]interface{}{"allowed_origins": []string{"www.example.org"}, "debug": true},
	}
	defer config.Delete(owner)
	err := config.Save(owner)
	c.Check(err, IsNil)
	c.Check(DeletePluginsByService(service.Subdomain), IsNil)
}
