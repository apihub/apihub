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

func (s *S) SetUpSuite(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)
}

func (s *S) TearDownSuite(c *C) {
}

var team account.Team
var owner account.User
var alice account.User
var service account.Service

func (s *S) SetUpTest(c *C) {
	team = account.Team{Name: "Backstage Team", Alias: "backstage"}
	alice = account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	owner = account.User{Name: "Owner", Email: "owner@example.org", Password: "123456"}
	service = account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage", Transformers: []string{}}
}

// Run the tests in memory
func setUpMemoryTest(s *S) {
	s.store = mem.New()
	account.NewStorable = func() (account.Storable, error) {
		return s.store, nil
	}
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_account_test",
	}
	account.NewStorable = func() (account.Storable, error) {
		return mongore.New(cfg)
	}
}
