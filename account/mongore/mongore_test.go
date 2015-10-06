package mongore

import (
	"testing"

	"github.com/apihub/apihub/account/test"
	. "gopkg.in/check.v1"
)

func TestMongore(t *testing.T) {
	config := Config{
		Host:         "127.0.0.1:27017",
		DatabaseName: "apihub_mongore_test",
	}

	m := New(config)
	Suite(&test.StorableSuite{Storage: m})
	TestingT(t)

	m.openSession().Collection("db").Database.DropDatabase()
}
