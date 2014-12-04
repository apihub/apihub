package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/albertoleal/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()
	defer account.DeleteGroupByName("Group")

	payload := `{"name": "Group"}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	s.env[CurrentUser] = user
	response := groupsController.CreateTeam(&web.C{Env: s.env}, s.recorder, req)
	c.Assert(response.StatusCode, Equals, 201)
	c.Assert(response.Payload, Matches, "^{\"id\":\".*?\",\"name\":\"Group\",\"users\":\\[\"alice\"\\],\"owner\":\"alice\"}$")
}

func (s *S) TestCreateTeamWhenUserIsNotSignedIn(c *C) {
	payload := `{"name": "Group"}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	response := groupsController.CreateTeam(&web.C{Env: s.env}, s.recorder, req)
	c.Assert(response.StatusCode, Equals, 400)
	c.Assert(response.Payload, Equals, "User is not signed in.")
}

func (s *S) TestCreateTeamWithInvalidPayloadFormat(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	payload := `"name": "Group"`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = user
	c.Assert(err, IsNil)
	webC := web.C{Env: s.env}
	groupsController.CreateTeam(&webC, s.recorder, req)
	expected := `{"status_code":400,"payload":"The request was bad-formed."}`
	key, _ := GetRequestError(&webC)
	body, _ := json.Marshal(key)
	c.Assert(string(body), Equals, expected)
}

func (s *S) TestDeleteTeam(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer owner.Delete()
	defer group.Delete()

	app := &Application{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id", app.Route(gg, "DeleteTeam"))

	g, _ := account.FindGroupByName(group.Name)
	req, err := http.NewRequest("DELETE", "/api/teams/"+g.Id.Hex(), nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = owner
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	expected := `{"name":"Group","users":["alice"],"owner":"alice"}`
	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Equals, expected)
}

func (s *S) TestDeleteTeamWhenUserIsNotOwner(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer owner.Delete()
	defer group.Delete()
	defer bob.Delete()

	app := &Application{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id", app.Route(gg, "DeleteTeam"))

	g, _ := account.FindGroupByName(group.Name)
	req, err := http.NewRequest("DELETE", "/api/teams/"+g.Id.Hex(), nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 403)
}

func (s *S) TestDeleteTeamIsNotFound(c *C) {
}
