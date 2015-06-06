package mongore

import (
	"testing"

	"github.com/backstage/backstage/account_new/test"
	. "gopkg.in/check.v1"
)

func TestMongore(t *testing.T) {
	config := Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "backstage_mongore_test",
	}

	m, _ := New(config)
	Suite(&test.StorableSuite{Storage: m})
	TestingT(t)

	m.(*Mongore).store.Collection("db").Database.DropDatabase()
	m.Close()
}
