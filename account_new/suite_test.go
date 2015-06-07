package account_new_test

import (
	"testing"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/account_new/mem"
	"github.com/backstage/backstage/account_new/mongore"
	. "gopkg.in/check.v1"
)

type S struct {
	store account_new.Storable
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

var team account_new.Team
var owner account_new.User
var alice account_new.User

func (s *S) SetUpTest(c *C) {
	team = account_new.Team{Name: "Backstage Team", Alias: "backstage"}
	alice = account_new.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	owner = account_new.User{Name: "Owner", Email: "owner@example.org", Password: "123456"}
}

// Run the tests in memory
func setUpMemoryTest(s *S) {
	s.store = mem.New()
	account_new.NewStorable = func() (account_new.Storable, error) {
		return s.store, nil
	}
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_account_test",
	}
	account_new.NewStorable = func() (account_new.Storable, error) {
		return mongore.New(cfg)
	}
}
