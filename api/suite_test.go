package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/db"
	"github.com/tsuru/config"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

var teamsHandler *TeamsHandler
var usersHandler *UsersHandler

var alice *account.User
var bob *account.User
var mary *account.User
var owner *account.User
var team *account.Team

func Test(t *testing.T) { TestingT(t) }

type S struct {
	Api *Api
	env        map[string]interface{}
	handler    http.HandlerFunc
	recorder   *httptest.ResponseRecorder
	router     *web.Mux
}

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_api_test")
}

func (s *S) SetUpTest(c *C) {
	s.Api = &Api{}
	teamsHandler = &TeamsHandler{}
	usersHandler = &UsersHandler{}

	s.recorder = httptest.NewRecorder()
	s.env = map[string]interface{}{}
	s.router = web.New()
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	s.router.Use(ErrorMiddleware)

	alice = &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	bob = &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	mary = &account.User{Name: "Bob", Email: "mary@example.org", Username: "mary", Password: "123456"}
	owner = &account.User{Name: "Bob", Email: "owner@example.org", Username: "owner", Password: "123456"}
	team = &account.Team{Name: "Team"}
}

func (s *S) TearDownSuite(c *C) {
	storage, err := db.Conn()
	c.Assert(err, IsNil)
	defer storage.Close()
	config.Unset("database:url")
	config.Unset("database:name")
	s.router.Abandon(ErrorMiddleware)
}

func (s *S) signIn() {

}

var _ = Suite(&S{})
