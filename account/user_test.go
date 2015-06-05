package account_test

import (
	"encoding/json"

	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

var alice = User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}

func (s *S) TestCreateUser(c *C) {
	defer alice.Delete()
	err := alice.Save()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := User{}
	err := user.Save()
	c.Assert(err, Equals, errors.ErrUserMissingRequiredFields)
}

func (s *S) TestCreateUserShouldMaskThePassword(c *C) {
	defer alice.Delete()
	alice.Save()

	foundUser, _ := FindUserByEmail(alice.Email)
	c.Assert(foundUser.Password, Not(Equals), "123456")
}

func (s *S) TestChangePassword(c *C) {
	defer alice.Delete()
	alice.Save()
	p1 := alice.Password
	alice.Password = "654321"
	err := alice.ChangePassword()
	c.Assert(err, IsNil)
	p2 := alice.Password
	c.Assert(p1, Not(Equals), p2)
}

func (s *S) TestExists(c *C) {
	defer alice.Delete()
	alice.Save()

	valid := alice.Exists()
	c.Assert(valid, Equals, true)
}

func (s *S) TestExistsWhenUserDoesNotExistInTheDB(c *C) {
	valid := alice.Exists()
	c.Assert(valid, Equals, false)
}

func (s *S) FindUserByEmail(c *C) {
	defer alice.Delete()
	alice.Save()

	foundUser, err := FindUserByEmail(alice.Email)
	c.Assert(err, IsNil)
	c.Assert(foundUser, NotNil)
}

func (s *S) TestFindUserWithInvalidUsername(c *C) {
	_, err := FindUserByEmail("bob@example.org")
	c.Assert(err, Equals, errors.ErrUserNotFound)
}

func (s *S) TestGetTeams(c *C) {
	defer alice.Delete()
	alice.Save()
	team := &Team{Name: "Team"}
	team.Save(&alice)
	defer DeleteTeamByAlias(team.Alias, &alice)
	g, err := alice.GetTeams()
	c.Assert(err, IsNil)
	c.Assert(g[0].Name, Equals, "Team")
}

func (s *S) TestGetServices(c *C) {
	defer alice.Delete()
	alice.Save()
	team := &Team{Name: "Team"}
	team.Save(&alice)
	defer DeleteTeamByAlias(team.Alias, &alice)
	service := &Service{Endpoint: "http://example.org/api", Subdomain: "_get_services", Transformers: []string{"cors"}}
	service.Save(&alice, team)
	defer DeleteServiceBySubdomain(service.Subdomain)
	g, err := alice.GetServices()
	cj, _ := json.Marshal(service)
	exp, _ := json.Marshal(g[0])
	c.Assert(err, IsNil)
	c.Assert(string(cj), Equals, string(exp))
}

func (s *S) TestUserToString(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	str := user.ToString()
	c.Assert(str, Equals, `{"name":"Alice","email":"alice@example.org","username":"alice"}`)
}
