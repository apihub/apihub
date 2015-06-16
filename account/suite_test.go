package account_test

import (
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/account/mem"
	"github.com/backstage/backstage/account/mongore"
	. "gopkg.in/check.v1"
)

type S struct {
	store account.Storable
}

var _ = Suite(&S{})

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *S) TearDownSuite(c *C) {
}

var app account.App
var alice account.User
var owner account.User
var service account.Service
var team account.Team

func (s *S) SetUpTest(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)

	team = account.Team{Name: "Backstage Team", Alias: "backstage"}
	alice = account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	owner = account.User{Name: "Owner", Email: "owner@example.org", Password: "123456"}
	service = account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage", Transformers: []string{}}
	app = account.App{ClientId: "ios", ClientSecret: "secret", Name: "Ios App", Team: team.Alias, Owner: owner.Email, RedirectUris: []string{"http://www.example.org/auth"}}
}

// Run the tests in memory
func setUpMemoryTest(s *S) {
	account.Storage(mem.New())
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	account.Storage(mongore.New(mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_account_test",
	}))
}
