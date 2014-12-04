package api

import (
	"net/http/httptest"
	"net/http"
	"testing"

	"github.com/albertoleal/backstage/db"
	"github.com/zenazn/goji/web"
	"github.com/tsuru/config"
	. "gopkg.in/check.v1"
)

var groupsController *GroupsController

func Test(t *testing.T) { TestingT(t) }

type S struct {
	recorder   *httptest.ResponseRecorder
	env        map[string]interface{}
	controller *UsersController //Should remove this.
	router   *web.Mux
	handler  http.HandlerFunc
}

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_api_test")
}

func (s *S) SetUpTest(c *C) {
	s.controller = &UsersController{}
	groupsController = &GroupsController{}
	s.recorder = httptest.NewRecorder()
	s.env = map[string]interface{}{}
	s.router = web.New()
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func (s *S) TearDownSuite(c *C) {
	storage, err := db.Conn()
	c.Assert(err, IsNil)
	defer storage.Close()
	config.Unset("database:url")
	config.Unset("database:name")
}

var _ = Suite(&S{})