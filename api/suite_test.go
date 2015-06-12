package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/account/mem"
	"github.com/backstage/backstage/account/mongore"
	"github.com/backstage/backstage/api"
	. "gopkg.in/check.v1"
)

var httpClient HTTPClient

func Test(t *testing.T) { TestingT(t) }

var service account.Service
var team account.Team
var user account.User

type S struct {
	api        *api.Api
	authHeader string
	store      func() (account.Storable, error)
	server     *httptest.Server
}

func (s *S) SetUpSuite(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)

	s.api = api.NewApi(s.store)
	s.server = httptest.NewServer(s.api.Handler())
	httpClient = NewHTTPClient(s.server.URL)

	team = account.Team{Name: "Backstage Team", Alias: "backstage"}
	service = account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage"}
}

func (s *S) SetUpTest(c *C) {
	user = account.User{Name: "Bob", Email: "bob@bar.example.org", Password: "secret"}
	user.Create()
	token, err := s.api.Login(user.Email, "secret")
	if err != nil {
		panic(err)
	}
	s.authHeader = fmt.Sprintf("%s %s", token.Type, token.Token)
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
	mem := mem.New()
	s.store = func() (account.Storable, error) {
		return mem, nil
	}
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_api_test",
	}
	s.store = func() (account.Storable, error) {
		return mongore.New(cfg)
	}
}

func testWithoutSignIn(reqArgs RequestArgs, c *C) {
	headers, code, body, err := httpClient.MakeRequest(reqArgs)

	c.Assert(string(body), Equals, `{"error":"unauthorized_access","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(code, Equals, http.StatusUnauthorized)
	c.Check(err, IsNil)
}
