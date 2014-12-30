package account

import (
	"testing"

	"github.com/tsuru/config"
	. "gopkg.in/check.v1"

	"github.com/backstage/backstage/db"
)

var (
	team  *Team
	owner *User
	service *Service
	client *Client
)

type S struct{}

var _ = Suite(&S{})

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *S) SetUpSuite(c *C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "backstage_db_test")
}

func (s *S) TearDownSuite(c *C) {
	conn, err := db.Conn()
	c.Assert(err, IsNil)
	defer conn.Close()
	config.Unset("database:url")
	config.Unset("database:name")
	// conn.Collection("services").Database.DropDatabase()
}

func (s *S) SetUpTest(c *C) {
	team = &Team{Name: "Team", Alias: "Alias"}
	owner = &User{Name: "Owner", Username: "owner", Email: "owner@example.org", Password: "123456"}
	service = &Service{Endpoint: "http://example.org/api", Subdomain: "backstage"}
	client = &Client{Id: "backstage", Secret: "SuperSecret", Name: "Backstage", RedirectUri: "http://example.org/auth"}
}
