package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RangelReale/osin"
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/log"
	"github.com/tsuru/config"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

var oAuthHandler *OAuthHandler
var servicesHandler *ServicesHandler
var teamsHandler *TeamsHandler
var usersHandler *UsersHandler
var clientsHandler *ClientsHandler
var pluginsHandler *PluginsHandler

var alice *account.User
var bob *account.User
var mary *account.User
var owner *account.User
var service *account.Service
var client *account.Client
var team *account.Team

var osinClient *osin.DefaultClient = &osin.DefaultClient{
	Id:          "test-1234",
	Secret:      "super-secret-string",
	RedirectUri: "http://www.example.org/auth",
}

var authorizeData *osin.AuthorizeData = &osin.AuthorizeData{
	Client:      osinClient,
	Code:        "test-123456789",
	ExpiresIn:   3600,
	CreatedAt:   bson.Now(),
	RedirectUri: "http://www.example.org/auth",
}

var accessData *osin.AccessData = &osin.AccessData{
	Client:        osinClient,
	AuthorizeData: authorizeData,
	AccessToken:   "test-123456",
	RefreshToken:  "test-refresh-7890",
	ExpiresIn:     3600,
	CreatedAt:     bson.Now(),
}

func Test(t *testing.T) { TestingT(t) }

type S struct {
	Api          *Api
	env          map[string]interface{}
	handler      http.HandlerFunc
	recorder     *httptest.ResponseRecorder
	router       *web.Mux
	oAuthStorage *OAuthMongoStorage
}

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_api_test")
	log.Logger.Disable()
}

func (s *S) SetUpTest(c *C) {
	s.Api = &Api{Config: &Config{}}
	teamsHandler = &TeamsHandler{}
	usersHandler = &UsersHandler{}
	servicesHandler = &ServicesHandler{}
	clientsHandler = &ClientsHandler{}
	oAuthHandler = &OAuthHandler{}
	pluginsHandler = &PluginsHandler{}

	s.recorder = httptest.NewRecorder()
	s.env = map[string]interface{}{}

	s.router = web.New()
	s.router.Post("/api/clients", s.Api.route(clientsHandler, "CreateClient"))
	s.router.Put("/api/clients/:id", s.Api.route(clientsHandler, "UpdateClient"))
	s.router.Get("/api/clients/:id", s.Api.route(clientsHandler, "GetClientInfo"))
	s.router.Delete("/api/clients/:id", s.Api.route(clientsHandler, "DeleteClient"))

	s.router.Post("/api/services", s.Api.route(servicesHandler, "CreateService"))
	s.router.Get("/api/services", s.Api.route(servicesHandler, "GetUserServices"))
	s.router.Put("/api/services/:subdomain", s.Api.route(servicesHandler, "UpdateService"))
	s.router.Get("/api/services/:subdomain", s.Api.route(servicesHandler, "GetServiceInfo"))
	s.router.Delete("/api/services/:subdomain", s.Api.route(servicesHandler, "DeleteService"))

	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	s.oAuthStorage = &OAuthMongoStorage{}

	alice = &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	bob = &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	mary = &account.User{Name: "Mary", Email: "mary@example.org", Username: "mary", Password: "123456"}
	owner = &account.User{Name: "Owner", Email: "owner@example.org", Username: "owner", Password: "123456"}
	team = &account.Team{Name: "Team", Alias: "team"}
	service = &account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage"}
	client = &account.Client{Id: "backstage", Secret: "SuperSecret", Name: "Backstage", RedirectUri: "http://example.org/auth"}
}

func (s *S) TearDownSuite(c *C) {
	storage, err := db.Conn()
	c.Assert(err, IsNil)
	defer storage.Close()
	config.Unset("database:url")
	config.Unset("database:name")
}

var _ = Suite(&S{})
