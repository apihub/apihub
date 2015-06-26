package db

import (
	"testing"

	"github.com/tsuru/config"
	. "gopkg.in/check.v1"
)

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type S struct {
	etcd *Etcd
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_db_test")
	s.etcd, _ = NewEtcd("/db_test", &EtcdConfig{Machines: []string{"http://127.0.0.1:2379"}})
}

func (s *S) TearDownSuite(c *C) {
	storage, err := Conn()
	c.Assert(err, IsNil)
	defer storage.Close()
	config.Unset("database:url")
	config.Unset("database:name")
	s.etcd.Close()
}
