package auth_test

import (
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/account/mem"
	"github.com/backstage/backstage/account/mongore"
	"github.com/backstage/backstage/auth"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	store func() (account.Storable, error)
	auth  auth.Authenticatable
}

func (s *S) SetUpSuite(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)
}

func (s *S) SetUpTest(c *C) {
	s.auth = auth.NewAuth(s.store)
}

var _ = Suite(&S{})

// Run the tests in memory
func setUpMemoryTest(s *S) {
	store := mem.New()
	s.store = func() (account.Storable, error) {
		return store, nil
	}
	account.NewStorable = s.store
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_auth_test",
	}
	s.store = func() (account.Storable, error) {
		return mongore.New(cfg)
	}
	account.NewStorable = s.store
}
