package db

import (
	"github.com/apihub/apihub/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestEtcdBasicOperationsKey(c *C) {
	err := s.etcd.SetKey("services", `{"name": "Alice"}`, 0)
	c.Check(err, IsNil)

	k, _ := s.etcd.GetKey("services")
	c.Assert(k, Equals, `{"name": "Alice"}`)

	err = s.etcd.DeleteKey("services")
	c.Check(err, IsNil)

	k, err = s.etcd.GetKey("services")
	c.Assert(k, Equals, "")
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}
