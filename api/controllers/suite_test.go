package controllers

import (
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/albertoleal/backstage/db"
	"github.com/tsuru/config"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	recorder   *httptest.ResponseRecorder
	env        map[string]interface{}
	controller *UsersController
}

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_api_controllers_test")
}

func (s *S) TearDownSuite(c *C) {
	storage, err := db.Conn()
	c.Assert(err, IsNil)
	defer storage.Close()
	config.Unset("database:url")
	config.Unset("database:name")
}

var _ = Suite(&S{})
