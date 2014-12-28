package auth

import (
	"github.com/backstage/backstage/db"
	"github.com/tsuru/config"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_auth_test")
}

func (s *S) TearDownSuite(c *C) {
	storage, err := db.Conn()
	c.Assert(err, IsNil)
	defer storage.Close()
	config.Unset("database:url")
	config.Unset("database:name")
}

var _ = Suite(&S{})
