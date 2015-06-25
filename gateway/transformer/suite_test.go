package transformer

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	transformers Transformers
}

func (s *S) SetUpTest(c *C) {
	s.transformers = make(map[string]Transformer)
}

var _ = Suite(&S{})
