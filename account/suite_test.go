package account_test

import (
	"testing"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/account/mem"
	"github.com/backstage/maestro/account/mongore"
	. "github.com/backstage/maestro/log"
	. "gopkg.in/check.v1"
)

type S struct {
	store account.Storable
}

var _ = Suite(&S{})

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *S) SetUpSuite(c *C) {
	Logger.Disable()
}

var app account.App
var alice account.User
var owner account.User
var pluginConfig account.Plugin
var service account.Service
var team account.Team
var hook account.Hook

func (s *S) SetUpTest(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)

	team = account.Team{Name: "Backstage Team", Alias: "backstage", Services: []account.Service{}, Apps: []account.App{}}
	alice = account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	owner = account.User{Name: "Owner", Email: "owner@example.org", Password: "123456"}
	service = account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage", Transformers: []string{}}
	app = account.App{ClientId: "ios", ClientSecret: "secret", Name: "Ios App", Team: team.Alias, Owner: owner.Email, RedirectUris: []string{"http://www.example.org/auth"}}
	pluginConfig = account.Plugin{Name: "Plugin Config", Service: service.Subdomain, Config: make(map[string]interface{})}
	hook = account.Hook{Name: "service.update", Events: []string{"service.update"}, Config: account.HookConfig{Address: "http://www.example.org"}}
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
