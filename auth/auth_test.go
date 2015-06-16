package auth_test

import (
	"testing"

	"github.com/backstage/apimanager/account"
	"github.com/backstage/apimanager/account/mem"
	"github.com/backstage/apimanager/account/mongore"
	"github.com/backstage/apimanager/auth"
	"github.com/backstage/apimanager/auth/test"
	. "gopkg.in/check.v1"
)

func TestAuth(t *testing.T) {
	// store := setUpMemoryTest()
	store := setUpMongoreTest()
	auth := auth.NewAuth(store)

	Suite(&test.AuthenticatableSuite{Auth: auth})
	TestingT(t)
}

// // Run the tests in memory
func setUpMemoryTest() account.Storable {
	store := mem.New()
	account.Storage(store)
	return store
}

// Run the tests using MongoRe
func setUpMongoreTest() account.Storable {
	cfg := mongore.Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_auth_test",
	}
	store := mongore.New(cfg)
	account.Storage(store)
	return store
}
