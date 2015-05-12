package gateway

import (
	"fmt"
	"testing"

	"github.com/backstage/backstage/db"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	Settings *Settings
}

func (s *S) SetUpTest(c *C) {
	s.Settings = &Settings{
		Host:        "test.backstage.dev",
		Port:        ":4567",
		ChannelName: "services",
	}
}

func (s *S) AddToken(token string, expires int, data map[string]interface{}) {
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	conn.Tokens(token, expires, data)
}

func (s *S) DeleteToken(token string) {
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	conn.DeleteToken(token)
}

var _ = Suite(&S{})
