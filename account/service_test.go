package account

import (
	. "gopkg.in/check.v1"
	"github.com/backstage/backstage/errors"
)

func (s *S) TestCreateServiceNewService(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := Service{
		Endpoint:  "http://example.org/api",
		Subdomain: "BACKSTAGE",
	}
	err := service.Save(owner, team)
	defer service.Delete()

	c.Check(service.Subdomain, Equals, "backstage")
	_, ok := err.(*errors.ValidationError)
	c.Check(ok, Equals, false)
}

func (s *S) TestCannotCreateServiceWhenSubdomainAlreadyExists(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := Service{
		Endpoint:  "http://example.org/api",
		Subdomain: "backstage",
	}
	err := service.Save(owner, team)
	defer service.Delete()
	c.Check(err, IsNil)

	service2 := Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}
	err = service2.Save(owner, team)
	c.Check(err, NotNil)

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	c.Assert(e.Payload, Equals, "There is another service with this subdomain.")
}

func (s *S) TestCannotCreateServiceAServiceWithoutRequiredFields(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := &Service{Subdomain: "backstage"}
	err := service.Save(owner, team)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Endpoint cannot be empty.")

	service = &Service{}
	err = service.Save(owner, team)
	e = err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Subdomain cannot be empty.")
}

func (s *S) TestDeleteServiceANonExistingService(c *C) {
	service := &Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}
	err := service.Delete()

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	c.Assert(e.Payload, Equals, "Service not found.")
}

func (s *S) TestDeleteServiceAnExistingService(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := &Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}

	count, _ := CountService()
	c.Assert(count, Equals, 0)

	service.Save(owner, team)
	count, _ = CountService()
	c.Assert(count, Equals, 1)

	service.Delete()
	count, _ = CountService()
	c.Assert(count, Equals, 0)
}

func (s *S) TestFindServiceBySubdomainByAlias(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := &Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}

	defer service.Delete()
	service.Save(owner, team)
	se, _ := FindServiceBySubdomain(service.Subdomain)
	c.Assert(se.Subdomain, Equals, service.Subdomain)
}

func (s *S) TestFindServiceBySubdomainWithInvalidName(c *C) {
	_, err := FindServiceBySubdomain("Non Existing Service")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Service not found.")
}

func (s *S) TestDeleteServiceBySubdomain(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := &Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}
	defer service.Delete()
	err := service.Save(owner, team)
	c.Assert(err, IsNil)
	err = DeleteServiceBySubdomain(service.Subdomain)
	c.Assert(err, IsNil)
}

func (s *S) TestDeleteServiceBySubdomainWithInvalidSubdomain(c *C) {
	err := DeleteServiceBySubdomain("Non existing service")
	c.Assert(err, NotNil)
	e := err.(*errors.ValidationError)
	c.Assert(e.Payload, Equals, "Service not found.")
}

func (s *S) TestFindServicesByTeam(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := &Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}

	defer service.Delete()
	service.Save(owner, team)
	se, _ := FindServicesByTeam(team.Alias)
	c.Assert(len(se), Equals, 1)
	c.Assert(se[0].Subdomain, Equals, "backstage")
}

func (s *S) TestFindServicesByTeamWithoutElements(c *C) {
	se, _ := FindServicesByTeam("non-existing-team")
	c.Assert(len(se), Equals, 0)
}
