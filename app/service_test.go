package service

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

func (s *S) TestCreateNewService(c *C) {
	service := Service{Name: "Backstage",
		Endpoint:  map[string]interface{}{"latest": "http://example.org/api"},
		Subdomain: "BACKSTAGE",
	}
	err := Create(&service)
	defer Delete(&service)

	c.Check(service.Subdomain, Equals, "backstage")
	_, ok := err.(*errors.ValidationError)
	c.Check(ok, Equals, false)
}

func (s *S) TestCannotCreateServiceWhenSubdomainAlreadyExists(c *C) {
	service := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}
	err := Create(&service)
	defer Delete(&service)
	c.Check(err, IsNil)

	service2 := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}
	err = Create(&service2)
	c.Check(err, NotNil)

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	msg := "There is another service with this subdomain."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestCannotCreateAServiceWithoutRequiredFields(c *C) {
	service := Service{Subdomain: "backstage"}
	err := Create(&service)
	e := err.(*errors.ValidationError)
	msg := "Endpoint cannot be empty."
	c.Assert(e.Message, Equals, msg)

	service = Service{}
	err = Create(&service)
	e = err.(*errors.ValidationError)
	msg = "Subdomain cannot be empty."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestDeleteANonExistingService(c *C) {
	service := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}
	err := Delete(&service)

	e, ok := err.(*errors.ValidationError)
	c.Assert(ok, Equals, true)
	msg := "Document not found."
	c.Assert(e.Message, Equals, msg)
}

func (s *S) TestDeleteAnExistingService(c *C) {
	service := Service{Subdomain: "backstage",
		Endpoint: map[string]interface{}{"latest": "http://example.org/api"},
	}

	count, _ := Count()
	c.Assert(count, Equals, 0)

	Create(&service)
	count, _ = Count()
	c.Assert(count, Equals, 1)

	Delete(&service)
	count, _ = Count()
	c.Assert(count, Equals, 0)
}
