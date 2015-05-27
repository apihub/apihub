package mongore

import (
	"testing"

	"github.com/backstage/backstage/account/test"
	. "gopkg.in/check.v1"
)

func TestMongore(t *testing.T) {
	Suite(&test.StorableSuite{Storage: New()})
	TestingT(t)
}
