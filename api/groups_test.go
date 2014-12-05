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

	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id", api.Route(gg, "DeleteTeam"))

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

	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id", api.Route(gg, "DeleteTeam"))

	g, _ := account.FindGroupByName(group.Name)
	req, err := http.NewRequest("DELETE", "/api/teams/"+g.Id.Hex(), nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, "Team not found or you're not the owner.")
}

func (s *S) TestDeleteTeamIsNotFound(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id", api.Route(gg, "DeleteTeam"))
	req, err := http.NewRequest("DELETE", "/api/teams/invalid-id", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, "Team not found or you're not the owner.")
}

func (s *S) TestGetUserTeams(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(bob)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Get("/api/teams", api.Route(gg, "GetUserTeams"))
	req, err := http.NewRequest("GET", "/api/teams", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Matches, "^\\[{\"id\":\".*?\",\"name\":\"Group\",\"users\":\\[\"bob\"\\],\"owner\":\"bob\"}\\]$")
}

func (s *S) TestGetUserTeamsWhenUserIsNotSignedIn(c *C) {
	api := &Api{}
	gg := &GroupsController{}
	s.router.Get("/api/teams", api.Route(gg, "GetUserTeams"))
	req, err := http.NewRequest("GET", "/api/teams", nil)
	c.Assert(err, IsNil)
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "User is not signed in.")
}

func (s *S) TestGetTeamInfo(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(bob)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Get("/api/teams/:id", api.Route(gg, "GetTeamInfo"))
	g, _ := account.FindGroupByName(group.Name)
	req, err := http.NewRequest("GET", "/api/teams/"+g.Id.Hex(), nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Group\",\"users\":\\[\"bob\"\\],\"owner\":\"bob\"}$")
}

func (s *S) TestGetTeamInfoWhenTeamNotFound(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(bob)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Get("/api/teams/:id", api.Route(gg, "GetTeamInfo"))
	req, err := http.NewRequest("GET", "/api/teams/invalid-id", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "Team not found.")
}

func (s *S) TestGetTeamInfoWhenIsNotMemberOfTheTeam(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Get("/api/teams/:id", api.Route(gg, "GetTeamInfo"))
	g, _ := account.FindGroupByName(group.Name)
	req, err := http.NewRequest("GET", "/api/teams/"+g.Id.Hex(), nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "You do not belong to this team!")
}

func (s *S) TestTeamInfoWhenUserIsNotSignedIn(c *C) {
	api := &Api{}
	gg := &GroupsController{}
	s.router.Get("/api/teams/:id", api.Route(gg, "GetTeamInfo"))
	req, err := http.NewRequest("GET", "/api/teams/1", nil)
	c.Assert(err, IsNil)
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "User is not signed in.")
}

func (s *S) TestAddUsersToTeam(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer account.DeleteGroupByName(group.Name)

	api := &Api{}
	gg := &GroupsController{}
	s.router.Post("/api/teams/:id/users", api.Route(gg, "AddUsersToTeam"))
	g, _ := account.FindGroupByName(group.Name)
	payload := `{"users": ["bob"]}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams/"+g.Id.Hex()+"/users", b)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = owner
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 201)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Group\",\"users\":\\[\"alice\",\"bob\"\\],\"owner\":\"alice\"}$")
}

func (s *S) TestAddUserToTeamWithInvalidPaylod(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer account.DeleteGroupByName(group.Name)

	api := &Api{}
	gg := &GroupsController{}
	s.router.Post("/api/teams/:id/users", api.Route(gg, "AddUsersToTeam"))
	g, _ := account.FindGroupByName(group.Name)
	payload := `{"members": ["bob"]}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams/"+g.Id.Hex()+"/users", b)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = owner
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "The request was bad-formed.")
}

func (s *S) TestAddUserToTeamWhenTeamNotFound(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(bob)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Post("/api/teams/:id/users", api.Route(gg, "AddUsersToTeam"))
	req, err := http.NewRequest("POST", "/api/teams/invalid-id/users", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "Team not found.")
}

func (s *S) TestAddUserToTeamWhenUserDoesNotBelongToIt(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	g, _ := account.FindGroupByName(group.Name)
	req, err := http.NewRequest("POST", "/api/teams/"+g.Id.Hex()+"/users", nil)
	s.router.Post("/api/teams/:id/users", api.Route(gg, "AddUsersToTeam"))
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, "You do not belong to this team!")
}

func (s *S) TestAddUsersToTeamWhenUserIsNotSignedIn(c *C) {
	api := &Api{}
	gg := &GroupsController{}
	s.router.Post("/api/teams/:id/users", api.Route(gg, "AddUsersToTeam"))
	req, err := http.NewRequest("POST", "/api/teams/invalid-id/users", nil)
	c.Assert(err, IsNil)
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "User is not signed in.")
}

func (s *S) TestRemoveUsersFromTeam(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer account.DeleteGroupByName(group.Name)

	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id/users", api.Route(gg, "RemoveUsersFromTeam"))
	g, _ := account.FindGroupByName(group.Name)
	payload := `{"users": ["bob"]}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("DELETE", "/api/teams/"+g.Id.Hex()+"/users", b)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = owner
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Group\",\"users\":\\[\"alice\"\\],\"owner\":\"alice\"}$")
}

func (s *S) TestRemoveUsersFromTeamWithInvalidPaylod(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer account.DeleteGroupByName(group.Name)

	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id/users", api.Route(gg, "RemoveUsersFromTeam"))
	g, _ := account.FindGroupByName(group.Name)
	payload := `{"members": ["bob"]}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("DELETE", "/api/teams/"+g.Id.Hex()+"/users", b)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = owner
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "The request was bad-formed.")
}

func (s *S) TestRemoveUsersFromTeamWhenTeamNotFound(c *C) {
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(bob)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id/users", api.Route(gg, "RemoveUsersFromTeam"))
	req, err := http.NewRequest("DELETE", "/api/teams/invalid-id/users", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "Team not found.")
}

func (s *S) TestRemoveUsersFromTeamWhenUserDoesNotBelongToIt(c *C) {
	owner := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	owner.Save()
	defer owner.Delete()
	bob := &account.User{Name: "Bob", Email: "bob@example.org", Username: "bob", Password: "123456"}
	bob.Save()
	defer bob.Delete()
	group := &account.Group{Name: "Group"}
	group.Save(owner)
	defer group.Delete()

	api := &Api{}
	gg := &GroupsController{}
	g, _ := account.FindGroupByName(group.Name)
	s.router.Delete("/api/teams/:id/users", api.Route(gg, "RemoveUsersFromTeam"))
	req, err := http.NewRequest("DELETE", "/api/teams/"+g.Id.Hex()+"/users", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = bob
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, "You do not belong to this team!")
}

func (s *S) TestRemoveUserFromTeamWhenUserIsNotSignedIn(c *C) {
	api := &Api{}
	gg := &GroupsController{}
	s.router.Delete("/api/teams/:id/users", api.Route(gg, "RemoveUsersFromTeam"))
	req, err := http.NewRequest("DELETE", "/api/teams/invalid-id/users", nil)
	c.Assert(err, IsNil)
	cc := web.C{Env: s.env}
	s.router.ServeHTTPC(cc, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, "User is not signed in.")
}
