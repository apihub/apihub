package account

import (
	"encoding/json"

	"github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	err := user.Save()
	c.Assert(err, IsNil)
}

func (s *S) TestCreateUserWithSameEmail(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	user2 := &User{Name: "Bob", Email: "alice@example.org", Username: "bob", Password: "123456"}
	err := user2.Save()
	e := err.(*errors.ValidationError)
	msg := "Someone already has that email/username. Could you try another?"
	c.Assert(e.Payload, Equals, msg)
}

func (s *S) TestCreateUserWithSameUsername(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	user2 := &User{Name: "Bob", Email: "bob@example.org", Username: "alice", Password: "123456"}
	err := user2.Save()
	e := err.(*errors.ValidationError)
	msg := "Someone already has that email/username. Could you try another?"
	c.Assert(e.Payload, Equals, msg)
}

func (s *S) TestCreateUserWithoutRequiredFields(c *C) {
	user := User{}
	err := user.Save()
	e := err.(*errors.ValidationError)
	msg := "Name/Email/Username/Password cannot be empty."
	c.Assert(e.Payload, Equals, msg)
}

func (s *S) TestCreateUserShouldMaskThePassword(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	foundUser, _ := FindUserByEmail(user.Email)
	c.Assert(foundUser.Password, Not(Equals), "123456")
}

func (s *S) TestValid(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	valid := user.Valid()
	c.Assert(valid, Equals, true)
}

func (s *S) TestValidWhenUserDoesNotExistInTheDB(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	valid := user.Valid()
	c.Assert(valid, Equals, false)
}

func (s *S) FindUserByEmail(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	foundUser, err := FindUserByEmail(user.Email)
	c.Assert(err, IsNil)
	c.Assert(foundUser, NotNil)
}

func (s *S) TestFindUserWithInvalidUsername(c *C) {
	user := User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()

	_, err := FindUserByEmail("bob@example.org")
	e := err.(*errors.ValidationError)
	msg := "User not found"
	c.Assert(e.Payload, Equals, msg)
}

func (s *S) TestGetTeams(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()
	team := &Team{Name: "Team"}
	team.Save(user)
	defer DeleteTeamByAlias(team.Alias, user)
	g, err := user.GetTeams()
	c.Assert(err, IsNil)
	c.Assert(g[0].Name, Equals, "Team")
}

func (s *S) TestGetServices(c *C) {
	user := &User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	defer user.Delete()
	user.Save()
	team := &Team{Name: "Team"}
	team.Save(user)
	defer DeleteTeamByAlias(team.Alias, user)
	service := &Service{Endpoint: "http://example.org/api", Subdomain: "_get_services", Transformers: []string{}}
	service.Save(user, team)
	defer DeleteServiceBySubdomain(service.Subdomain)
	g, err := user.GetServices()
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
