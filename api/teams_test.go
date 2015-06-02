package api

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	alice.Save()
	defer alice.Delete()
	defer account.DeleteTeamByName("Team")

	payload := `{"name": "Team"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"alice@example.org\"\\],\"owner\":\"alice@example.org\"}$")
}

func (s *S) TestCreateTeamWithAlias(c *C) {
	alice.Save()
	defer alice.Delete()
	defer account.DeleteTeamByName("Team")

	payload := `{"name": "Team", "alias": "my alias"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"name\":\"Team\",\"alias\":\"my-alias\",\"users\":\\[\"alice@example.org\"\\],\"owner\":\"alice@example.org\"}$")
}

func (s *S) TestCreateTeamWhenUserIsNotSignedIn(c *C) {
	payload := `{"name": "Team"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestCreateTeamWithInvalidPayloadFormat(c *C) {
	alice.Save()
	defer alice.Delete()
	payload := `"name": "Team"`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams", s.Api.route(teamsHandler, "CreateTeam"))
	req, _ := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestUpdateTeam(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	payload := `{"name": "New Name"}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Put("/api/teams/:alias", s.Api.route(teamsHandler, "UpdateTeam"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+g.Alias, b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"name":"New Name","alias":"team","users":["owner@example.org"],"owner":"owner@example.org"}`)
	c.Assert(s.recorder.Code, Equals, http.StatusOK)
}

func (s *S) TestUpdateTeamWhenUserIsNotOwner(c *C) {
	bob.Save()
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer bob.Delete()

	payload := `{"name": "New Name"}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Put("/api/teams/:alias", s.Api.route(teamsHandler, "UpdateTeam"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+g.Alias, b)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestUpdateTeamWhenNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	payload := `{"name": "New Name"}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:alias", s.Api.route(teamsHandler, "UpdateTeam"))
	req, _ := http.NewRequest("PUT", "/api/teams/invalid-id", b)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestDeleteTeam(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias", s.Api.route(teamsHandler, "DeleteTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"name":"Team","alias":"team","users":["owner@example.org"],"owner":"owner@example.org"}`)
}

func (s *S) TestDeleteTeamWhenUserIsNotOwner(c *C) {
	bob.Save()
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer bob.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias", s.Api.route(teamsHandler, "DeleteTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"Only the owner has permission to perform this operation."}`)
}

func (s *S) TestDeleteTeamIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Delete("/api/teams/:alias", s.Api.route(teamsHandler, "DeleteTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/invalid-id", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"Only the owner has permission to perform this operation."}`)
}

func (s *S) TestGetUserTeams(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	s.router.Get("/api/teams", s.Api.route(teamsHandler, "GetUserTeams"))
	req, _ := http.NewRequest("GET", "/api/teams", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"items":[{"name":"Team","alias":"team","users":["owner@example.org"],"owner":"owner@example.org"}],"item_count":1}`)
}

func (s *S) TestGetUserTeamsWhenUserIsNotSignedIn(c *C) {
	s.router.Get("/api/teams", s.Api.route(teamsHandler, "GetUserTeams"))
	req, _ := http.NewRequest("GET", "/api/teams", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestGetTeamInfo(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Get("/api/teams/:alias", s.Api.route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\"\\],\"owner\":\"owner@example.org\"}$")
}

func (s *S) TestGetTeamInfoIncludeServices(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	g, _ := account.FindTeamByName(team.Name)
	s.router.Get("/api/teams/:alias", s.Api.route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\"\\],\"owner\":\"owner@example.org\",\"services\":\\[.*?\\]}$")
}

func (s *S) TestGetTeamInfoWhenTeamNotFound(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	s.router.Get("/api/teams/:alias", s.Api.route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/invalid-id", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestGetTeamInfoWhenIsNotMemberOfTheTeam(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Get("/api/teams/:alias", s.Api.route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+g.Alias, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestTeamInfoWhenUserIsNotSignedIn(c *C) {
	s.router.Get("/api/teams/:alias", s.Api.route(teamsHandler, "GetTeamInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/1", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
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
	s.router.Post("/api/teams/:alias/users", s.Api.route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\",\"bob@example.org\"\\],\"owner\":\"owner@example.org\"}$")
}

func (s *S) TestAddUserToTeamWithInvalidPayload(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `"users": ["` + bob.Email + `"]`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Post("/api/teams/:alias/users", s.Api.route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestAddUserToTeamWhenTeamNotFound(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	s.router.Post("/api/teams/:alias/users", s.Api.route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/invalid-id/users", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestAddUsersToTeamWhenUserDoesNotBelongToIt(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Post("/api/teams/:alias/users", s.Api.route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/"+g.Alias+"/users", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestAddUsersToTeamWhenUserIsNotSignedIn(c *C) {
	s.router.Post("/api/teams/:alias/users", s.Api.route(teamsHandler, "AddUsersToTeam"))
	req, _ := http.NewRequest("POST", "/api/teams/invalid-id/users", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
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
	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Matches, "^{\"name\":\"Team\",\"alias\":\"team\",\"users\":\\[\"owner@example.org\"\\],\"owner\":\"owner@example.org\"}$")
}

func (s *S) TestRemoveUsersFromTeamWithInvalidPayload(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `"members": ["` + bob.Email + `"]`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestRemoveUsersFromTeamWithInvalidKey(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByName(team.Name)

	payload := `{"members": ["` + bob.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"name":"Team","alias":"team","users":["owner@example.org"],"owner":"owner@example.org"}`)
}

func (s *S) TestRemoveUsersFromTeamWhenTeamNotFound(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/invalid-id/users", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestRemoveUsersFromTeamWhenUserDoesNotBelongToIt(c *C) {
	owner.Save()
	bob.Save()
	team.Save(owner)
	defer owner.Delete()
	defer bob.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	payload := `{"users": ["` + bob.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestRemoveUsersFromTeamWhenUserIsOwner(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	payload := `{"users": ["` + owner.Email + `"]}`
	b := strings.NewReader(payload)

	g, _ := account.FindTeamByName(team.Name)
	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+g.Alias+"/users", b)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"It is not possible to remove the owner from the team."}`)
}

func (s *S) TestRemoveUserFromTeamWhenUserIsNotSignedIn(c *C) {
	s.router.Delete("/api/teams/:alias/users", s.Api.route(teamsHandler, "RemoveUsersFromTeam"))
	req, _ := http.NewRequest("DELETE", "/api/teams/invalid-id/users", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}
