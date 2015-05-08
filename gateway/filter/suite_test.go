package filter

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	filters Filters
}

func (s *S) SetUpTest(c *C) {
	s.filters = make(map[string]Filter)
}

var _ = Suite(&S{})
