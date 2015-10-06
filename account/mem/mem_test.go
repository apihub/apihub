package mem

import (
	"testing"

	"github.com/apihub/apihub/account/test"
	. "gopkg.in/check.v1"
)

func TestMem(t *testing.T) {
	Suite(&test.StorableSuite{Storage: New()})
	TestingT(t)
}
