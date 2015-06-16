package mem

import (
	"testing"

	"github.com/backstage/apimanager/account/test"
	. "gopkg.in/check.v1"
)

func TestMem(t *testing.T) {
	Suite(&test.StorableSuite{Storage: New()})
	TestingT(t)
}
