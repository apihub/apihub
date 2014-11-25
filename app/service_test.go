package app

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
	conn.Collection("services").Database.DropDatabase()
}

func (s *S) TestCreateServiceNewService(c *C) {
	service := Service{Name: "Backstage",
		Endpoint:  map[string]interface{}{"latest": "http://example.org/api"},
		Subdomain: "BACKSTAGE",
	}
	err := CreateService(&service)
	defer DeleteService(&service)

	c.Check(service.Subdomain, Equals, "backstage")
	_, ok := err.(*errors.ValidationError)
	c.Check(ok, Equals, false)
}

func (s *S) TestCannotCreateServiceServiceWhenSubdomainAlreadyExists(c *C) {
	service := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}
	err := CreateService(&service)
	defer DeleteService(&service)
	c.Check(err, IsNil)

	service2 := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}
	err = CreateService(&service2)
	c.Check(err, NotNil)

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	msg := "There is another service with this subdomain."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCannotCreateServiceAServiceWithoutRequiredFields(c *C) {
	service := Service{Subdomain: "backstage"}
	err := CreateService(&service)
	e := err.(*errors.ValidationError)
	msg := "Endpoint cannot be empty."
	c.Assert(e.Message, Equals, msg)

	service = Service{}
	err = CreateService(&service)
	e = err.(*errors.ValidationError)
	msg = "Subdomain cannot be empty."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestDeleteServiceANonExistingService(c *C) {
	service := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}
	err := DeleteService(&service)

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	msg := "Document not found."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestDeleteServiceAnExistingService(c *C) {
	service := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}

	count, _ := CountService()
	c.Assert(count, Equals, 0)

	CreateService(&service)
	count, _ = CountService()
	c.Assert(count, Equals, 1)

	DeleteService(&service)
	count, _ = CountService()
	c.Assert(count, Equals, 0)
}
