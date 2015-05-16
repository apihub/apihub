package account

import (
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateMiddlewareConfigNewService(c *C) {
	service := Service{
		Endpoint:  "http://example.org/api",
		Subdomain: "_test_middleware_config_new_service",
	}
	config := &MiddlewareConfig{
		Name:    "cors",
		Service: service.Subdomain,
		Config:  map[string]interface{}{"allow_origins": []string{"www"}, "debug": true},
	}
	defer config.Delete()
	err := config.Save()
	_, ok := err.(*errors.ValidationError)
	c.Check(ok, Equals, false)
}

func (s *S) TestCannotCreateMiddlewareConfigWithoutRequiredFields(c *C) {
	midd := &MiddlewareConfig{}
	err := midd.Save()
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Name cannot be empty.")

	midd = &MiddlewareConfig{Name: "foo"}
	err = midd.Save()
	e = err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Service cannot be empty.")
}

func (s *S) TestDeleteMiddConfigANonExistingMidd(c *C) {
	config := &MiddlewareConfig{
		Name:    "foo",
		Service: "backstage",
	}
	err := config.Delete()

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	c.Assert(e.Payload, Equals, "Middleware Config not found.")
}

// FIXME: ?
func (s *S) TestDeleteMiddCfgAnExisting(c *C) {
	config := &MiddlewareConfig{
		Name:    "cors",
		Service: "_test_midd",
		Config:  map[string]interface{}{"allow_origins": []string{"www"}, "debug": true},
	}
	err := config.Save()
	err = config.Delete()
	c.Check(err, IsNil)
}
