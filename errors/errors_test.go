package errors

import (
	"testing"

	. "gopkg.in/check.v1"
)

type S struct{}

var _ = Suite(&S{})

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *S) TestValidationError(c *C) {
	err := ValidationError{Message: "Something went wrong."}
	c.Assert(err.Error(), Equals, "Something went wrong.")
}
