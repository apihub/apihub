package api_new_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/account_new/mem"
	"github.com/backstage/backstage/account_new/mongore"
	"github.com/backstage/backstage/api_new"
	. "gopkg.in/check.v1"
)

var httpClient HTTPClient

func Test(t *testing.T) { TestingT(t) }

var user *account_new.User

type S struct {
	authHeader string
	store      func() (account_new.Storable, error)
	server     *httptest.Server
}

func (s *S) SetUpTest(c *C) {
	setUpMemoryTest(s)
	// setUpMongoreTest(s)

	api := api_new.NewApi(s.store)
	s.server = httptest.NewServer(api.Handler())
	httpClient = NewHTTPClient(s.server.URL)

	user = &account_new.User{Name: "Bob", Email: "bob@bar.example.org", Password: "secret"}
	user.Create()
	token, err := api.Login(user.Email, "secret")
	if err != nil {
		panic(err)
	}
	s.authHeader = fmt.Sprintf("%s %s", token.Type, token.Token)
}

func (s *S) TearDownTest(c *C) {
	s.server.Close()
	user.Delete()
}

var _ = Suite(&S{})

// Run the tests in memory
func setUpMemoryTest(s *S) {
	mem := mem.New()
	s.store = func() (account_new.Storable, error) {
		return mem, nil
	}
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_api_test",
	}
	s.store = func() (account_new.Storable, error) {
		return mongore.New(cfg)
	}
}
