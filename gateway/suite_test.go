package gateway

import (
	"fmt"
	"testing"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	Settings *Settings
}

var client *account.Client
var owner *account.User
var service *account.Service
var team *account.Team

func (s *S) SetUpTest(c *C) {
	s.Settings = &Settings{
		Host:        "test.backstage.dev",
		Port:        ":4567",
		ChannelName: "services",
	}
	client = &account.Client{Id: "backstage", Secret: "SuperSecret", Name: "Backstage", RedirectUri: "http://example.org/auth"}
	owner = &account.User{Name: "Owner", Email: "owner@example.org", Username: "owner", Password: "123456"}
	service = &account.Service{Endpoint: "http://example.org/api", Subdomain: "backstage"}
	team = &account.Team{Name: "Team", Alias: "team"}
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
