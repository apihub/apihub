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
	store account.Storable
	auth  auth.Authenticatable
}

func (s *S) SetUpTest(c *C) {
	// setUpMemoryTest(s)
	setUpMongoreTest(s)
	s.auth = auth.NewAuth(s.store)
}

var _ = Suite(&S{})

// Run the tests in memory
func setUpMemoryTest(s *S) {
	s.store = mem.New()
	account.Storage(s.store)
}

// Run the tests using MongoRe
func setUpMongoreTest(s *S) {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_auth_test",
	}
	s.store = mongore.New(cfg)
	account.Storage(s.store)
}
