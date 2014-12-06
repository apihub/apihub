package api

import (
	// "encoding/json"
	"net/http"
	"strings"

	"github.com/albertoleal/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	alice.Save()
	defer alice.Delete()
	defer account.DeleteTeamByName("Team")

	payload := `{"name": "Team"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.Route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 201)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"alice@example.org\"\\],\"owner\":\"alice@example.org\"}$")
}

func (s *S) TestCreateTeamWithAlias(c *C) {
	alice.Save()
	defer alice.Delete()
	defer account.DeleteTeamByName("Team")

	payload := `{"name": "Team", "alias": "my alias"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.Route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 201)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Team\",\"alias\":\"my-alias\",\"users\":\\[\"alice@example.org\"\\],\"owner\":\"alice@example.org\"}$")
}

func (s *S) TestCreateTeamWhenUserIsNotSignedIn(c *C) {
	payload := `{"name": "Team"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.Route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"User is not signed in."}`)
}

func (s *S) TestCreateTeamWithInvalidPayloadFormat(c *C) {
	alice.Save()
	defer alice.Delete()

	payload := `"name": "Team"`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.Route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"The request was bad-formed."}`)
}

func (s *S) TestDeleteTeam(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias", s.Api.Route(teamsHandler, "DeleteTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Equals, `{"name":"Team","alias":"team","users":["owner@example.org"],"owner":"owner@example.org"}`)
}

func (s *S) TestDeleteTeamWhenUserIsNotOwner(c *C) {
	bob.Save()
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()
	defer bob.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias", s.Api.Route(teamsHandler, "DeleteTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"Team not found or you're not the owner."}`)
}

func (s *S) TestDeleteTeamIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Delete("/api/teams/:alias", s.Api.Route(teamsHandler, "DeleteTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/invalid-id", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"Team not found or you're not the owner."}`)
}

func (s *S) TestGetUserTeams(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()

	s.router.Get("/api/teams", s.Api.Route(teamsHandler, "GetUserTeams"))
	req, _ := http.NewRequest("GET", "/api/teams", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Matches, "^\\[{\"id\":\".*?\",\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\"\\],\"owner\":\"owner@example.org\"}\\]$")
}

func (s *S) TestGetUserTeamsWhenUserIsNotSignedIn(c *C) {
	s.router.Get("/api/teams", s.Api.Route(teamsHandler, "GetUserTeams"))
	req, _ := http.NewRequest("GET", "/api/teams", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"User is not signed in."}`)
}

func (s *S) TestGetTeamInfo(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Get("/api/teams/:alias", s.Api.Route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\"\\],\"owner\":\"owner@example.org\"}$")
}

func (s *S) TestGetTeamInfoWhenTeamNotFound(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()

	s.router.Get("/api/teams/:alias", s.Api.Route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/invalid-id", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"Team not found."}`)
}

func (s *S) TestGetTeamInfoWhenIsNotMemberOfTheTeam(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer team.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Get("/api/teams/:alias", s.Api.Route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"You do not belong to this team!"}`)
}

func (s *S) TestTeamInfoWhenUserIsNotSignedIn(c *C) {
	s.router.Get("/api/teams/:alias", s.Api.Route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/1", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"User is not signed in."}`)
}

func (s *S) TestAddUsersToTeam(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer bob.Delete()
	defer owner.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `{"users": ["` + bob.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Post("/api/teams/:alias/users", s.Api.Route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 201)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\",\"bob@example.org\"\\],\"owner\":\"owner@example.org\"}$")
}

func (s *S) TestAddUserToTeamWithInvalidPaylod(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `{"members": ["` + bob.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Post("/api/teams/:alias/users", s.Api.Route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"The request was bad-formed."}`)
}

func (s *S) TestAddUserToTeamWhenTeamNotFound(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()

	s.router.Post("/api/teams/:alias/users", s.Api.Route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/invalid-id/users", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"Team not found."}`)
}

func (s *S) TestAddUsersToTeamWhenUserDoesNotBelongToIt(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer team.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Post("/api/teams/:alias/users", s.Api.Route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/"+g.Alias+"/users", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"You do not belong to this team!"}`)
}

func (s *S) TestAddUsersToTeamWhenUserIsNotSignedIn(c *C) {
	s.router.Post("/api/teams/:alias/users", s.Api.Route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/invalid-id/users", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"User is not signed in."}`)
}

func (s *S) TestRemoveUsersFromTeam(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer bob.Delete()
	defer owner.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `{"users": ["` + bob.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.Route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\".*?\",\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\"\\],\"owner\":\"owner@example.org\"}$")
}

func (s *S) TestRemoveUsersFromTeamWithInvalidPaylod(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `{"members": ["` + bob.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.Route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"The request was bad-formed."}`)
}

func (s *S) TestRemoveUsersFromTeamWhenTeamNotFound(c *C) {
	bob.Save()
	team.Save(bob)
	defer bob.Delete()
	defer team.Delete()

	s.router.Delete("/api/teams/:alias/users", s.Api.Route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/invalid-id/users", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"Team not found."}`)
}

func (s *S) TestRemoveUsersFromTeamWhenUserDoesNotBelongToIt(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer team.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.Route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"You do not belong to this team!"}`)
}

func (s *S) TestRemoveUsersFromTeamWhenUserIsOwner(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer team.Delete()

	payload := `{"users": ["` + owner.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.Route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"It is not possible to remove the owner from the team."}`)
}

func (s *S) TestRemoveUserFromTeamWhenUserIsNotSignedIn(c *C) {
	s.router.Delete("/api/teams/:alias/users", s.Api.Route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/invalid-id/users", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"User is not signed in."}`)
}
