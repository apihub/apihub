package account

import (
	"testing"

	"github.com/tsuru/config"
	. "gopkg.in/check.v1"

	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
)

type S struct{}

var _ = Suite(&S{})

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_db_test")
}

func (s *S) TearDownSuite(c *C) {
	conn, err := db.Conn()
	c.Assert(err, IsNil)
	defer conn.Close()
	config.Unset("database:url")
	config.Unset("database:name")
	// conn.Collection("services").Database.DropDatabase()
}

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

func (s *S) TestCannotCreateServiceServiceWhenSubdomainAlreadyExists(c *C) {
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
	message := "There is another service with this subdomain."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestCannotCreateServiceAServiceWithoutRequiredFields(c *C) {
	owner := &User{Email: "owner@example.org"}
	team := &Team{Name: "Team", Alias: "team"}
	service := &Service{Subdomain: "backstage"}
	err := service.Save(owner, team)
	e := err.(*errors.ValidationError)
	message := "Endpoint cannot be empty."
	c.Assert(e.Message, Equals, message)

	service = &Service{}
	err = service.Save(owner, team)
	e = err.(*errors.ValidationError)
	message = "Subdomain cannot be empty."
	c.Assert(e.Message, Equals, message)
}

func (s *S) TestDeleteServiceANonExistingService(c *C) {
	service := &Service{
		Subdomain: "backstage",
		Endpoint:  "http://example.org/api",
	}
	err := service.Delete()

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	message := "Document not found."
	c.Assert(e.Message, Equals, message)
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
	message := "Service not found."
	c.Assert(e.Message, Equals, message)
}
