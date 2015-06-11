package auth_new_test

import (
	"testing"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/account_new/mem"
	"github.com/backstage/backstage/account_new/mongore"
	"github.com/backstage/backstage/auth_new"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	store func() (account_new.Storable, error)
	auth  auth_new.Authenticatable
}

func (s *S) SetUpSuite(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)
}

func (s *S) SetUpTest(c *C) {
	s.auth = auth_new.NewAuth(s.store)
}

var _ = Suite(&S{})

// Run the tests in memory
func setUpMemoryTest(s *S) {
	store := mem.New()
	s.store = func() (account_new.Storable, error) {
		return store, nil
	}
	account_new.NewStorable = s.store
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_auth_test",
	}
	s.store = func() (account_new.Storable, error) {
		return mongore.New(cfg)
	}
	account_new.NewStorable = s.store
}
