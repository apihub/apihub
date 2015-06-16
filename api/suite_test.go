package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/backstage/apimanager/account"
	"github.com/backstage/apimanager/account/mem"
	"github.com/backstage/apimanager/account/mongore"
	"github.com/backstage/apimanager/api"
	. "gopkg.in/check.v1"
)

var httpClient HTTPClient

func Test(t *testing.T) { TestingT(t) }

var app account.App
var pluginConfig account.PluginConfig
var service account.Service
var team account.Team
var user account.User

type S struct {
	api        *api.Api
	authHeader string
	store      account.Storable
	server     *httptest.Server
}

// func (s *S) SetUpSuite(c *C) {
func (s *S) SetUpTest(c *C) {
	// setUpMongoreTest(s)
	setUpMemoryTest(s)

	s.api = api.NewApi(s.store)
	s.server = httptest.NewServer(s.api.Handler())
	httpClient = NewHTTPClient(s.server.URL)

	team = account.Team{Name: "Backstage Team", Alias: "backstage"}
	service = account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage"}
	user = account.User{Name: "Bob", Email: "bob@bar.example.org", Password: "secret"}
	app = account.App{ClientId: "ios", ClientSecret: "secret", Name: "Ios App", Team: team.Alias, Owner: user.Email, RedirectUris: []string{"http://www.example.org/auth"}}
	pluginConfig = account.PluginConfig{Name: "Plugin Config", Service: service.Subdomain, Config: map[string]interface{}{"version": 1}}

	user.Create()
	token, err := s.api.Login(user.Email, "secret")
	if err != nil {
		panic(err)
	}
	s.authHeader = fmt.Sprintf("%s %s", token.Type, token.AccessToken)
}

func (s *S) TearDownTest(c *C) {
	user.Delete()
}

func (s *S) TearDownSuite(c *C) {
	s.server.Close()
}

var _ = Suite(&S{})

// Run the tests in memory
func setUpMemoryTest(s *S) {
	s.store = mem.New()
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	s.store = mongore.New(mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_api_test",
	})
}

func testWithoutSignIn(reqArgs RequestArgs, c *C) {
	headers, code, body, err := httpClient.MakeRequest(reqArgs)

	c.Assert(string(body), Equals, `{"error":"unauthorized_access","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(code, Equals, http.StatusUnauthorized)
	c.Check(err, IsNil)
}
