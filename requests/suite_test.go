package requests_test

import (
	"testing"

	. "gopkg.in/check.v1"
)

type S struct {
}

var _ = Suite(&S{})

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
