package mem

import (
	"testing"

	"github.com/backstage/backstage/account/test"
	. "gopkg.in/check.v1"
)

func TestMem(t *testing.T) { TestingT(t) }

type MemSuite struct {
	suite test.StorageSuite
}

var _ = Suite(&MemSuite{})

func (s *MemSuite) SetUpTest(c *C) {
	storage := New()
	s.suite.Storage = storage
}

func (s *MemSuite) TestSaveToken(c *C) {
	s.suite.TestSaveToken(c)
}

func (s *MemSuite) TestGetToken(c *C) {
	s.suite.TestGetToken(c)
}

func (s *MemSuite) TestGetTokenWithNonExistingKey(c *C) {
	s.suite.TestGetTokenWithNonExistingKey(c)
}

func (s *MemSuite) TestDeleteToken(c *C) {
	s.suite.TestDeleteToken(c)
}

func (s *MemSuite) TestDeleteTokenWithNonExistingKey(c *C) {
	s.suite.TestDeleteTokenWithNonExistingKey(c)
}
