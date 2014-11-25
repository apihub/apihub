package storage

import (
	"testing"

	. "gopkg.in/check.v1"
)

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (ms *MySuite) TestOpenConnectsToDatabase(c *C) {
	storage, err := Open("127.0.0.1:27017", "backstage_db_test")
	c.Assert(err, IsNil)
	defer storage.session.Close()
}

func (ms *MySuite) TestOpenReuseSessionFromPoolForTheSameDatabase(c *C) {
	storage, err := Open("127.0.0.1:27017", "backstage_db_test")
	storage2, err := Open("127.0.0.1:27017", "backstage_db_test")
	c.Assert(err, IsNil)
	c.Assert(err, IsNil)
	defer storage.session.Close()
	defer storage2.session.Close()

	c.Assert(storage, Equals, storage2)
}

func (ms *MySuite) TestOpenDoNotReuseSessionFromPoolForDiffDatabase(c *C) {
	storage, err := Open("127.0.0.1:27017", "backstage_db_test")
	storage2, err := Open("127.0.0.1:27017", "backstage_db2_test")
	c.Assert(err, IsNil)
	c.Assert(err, IsNil)
	defer storage.session.Close()
	defer storage2.session.Close()

	c.Assert(storage, Not(Equals), storage2)
}
