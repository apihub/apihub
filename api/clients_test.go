package api

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateClient(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer account.DeleteClientByIdAndTeam("backstage", team.Alias)

	payload := `{"team": "` + team.Alias + `", "name": "Backstage", "redirect_uri": "http://www.example.org/auth"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/clients", s.Api.Route(clientsHandler, "CreateClient"))
	req, _ := http.NewRequest("POST", "/api/clients", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Matches, "^{\"id\":\"backstage\",\"secret\":\".*?\",\"name\":\"Backstage\",\"redirect_uri\":\"http://www.example.org/auth\",\"owner\":\"owner@example.org\",\"team\":\"team\"}$")
	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
}

func (s *S) TestCreateClientWithIdAndSecret(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer account.DeleteClientByIdAndTeam("awesome-id", team.Alias)

	payload := `{"team": "` + team.Alias + `", "id": "awesome id", "secret": "super secret", "name": "Backstage", "redirect_uri": "http://www.example.org/auth"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/clients", s.Api.Route(clientsHandler, "CreateClient"))
	req, _ := http.NewRequest("POST", "/api/clients", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"id":"awesome-id","secret":"super secret","name":"Backstage","redirect_uri":"http://www.example.org/auth","owner":"owner@example.org","team":"team"}`)
	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
}

func (s *S) TestCreateClientWhenUserIsNotSignedIn(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Post("/api/clients", s.Api.Route(clientsHandler, "CreateClient"))
	req, _ := http.NewRequest("POST", "/api/clients", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestCreateClientWhenTeamDoesNotExist(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	payload := `{"team": "invalid-team", "name": "Backstage", "redirect_uri": "http://www.example.org/auth"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/clients", s.Api.Route(clientsHandler, "CreateClient"))
	req, _ := http.NewRequest("POST", "/api/clients", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Team not found."}`)
}

func (s *S) TestCreateClientWithInvalidPayloadFormat(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	payload := `"name": "backstage"`
	b := strings.NewReader(payload)

	s.router.Post("/api/clients", s.Api.Route(clientsHandler, "CreateClient"))
	req, _ := http.NewRequest("POST", "/api/clients", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestDeleteClient(c *C) {
	owner.Save()
	team.Save(owner)
	client.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteClientByIdAndTeam(client.Id, team.Alias)
	defer owner.Delete()

	s.router.Delete("/api/clients/:id/:team", s.Api.Route(clientsHandler, "DeleteClient"))
	req, _ := http.NewRequest("DELETE", "/api/clients/"+client.Id+"/"+team.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"id":"backstage","secret":"SuperSecret","name":"Backstage","redirect_uri":"http://example.org/auth","owner":"owner@example.org","team":"team"}`)
}

func (s *S) TestDeleteClientWhenUserIsNotOwner(c *C) {
	alice.Save()
	owner.Save()
	team.Save(owner)
	client.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteClientByIdAndTeam(client.Id, team.Alias)
	defer alice.Delete()
	defer owner.Delete()

	s.router.Delete("/api/clients/:id/:team", s.Api.Route(clientsHandler, "DeleteClient"))
	req, _ := http.NewRequest("DELETE", "/api/clients/"+client.Id+"/"+team.Alias, nil)
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Client not found on this team."}`)
}

func (s *S) TestDeleteClientIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Delete("/api/clients/:id/:team", s.Api.Route(clientsHandler, "DeleteClient"))
	req, _ := http.NewRequest("DELETE", "/api/clients/invalid-client/"+team.Alias, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Client not found on this team."}`)
}

func (s *S) TestDeleteClientWhenTeamIsNotFound(c *C) {
	owner.Save()
	team.Save(owner)
	client.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteClientByIdAndTeam(client.Id, team.Alias)
	defer owner.Delete()

	s.router.Delete("/api/clients/:id/:team", s.Api.Route(clientsHandler, "DeleteClient"))
	req, _ := http.NewRequest("DELETE", "/api/clients/"+client.Id+"/invalid-client", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Client not found on this team."}`)
}

func (s *S) TestGetClientInfo(c *C) {
	owner.Save()
	team.Save(owner)
	client.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteClientByIdAndTeam(client.Id, team.Alias)
	defer owner.Delete()

	s.router.Get("/api/clients/:id/:team", s.Api.Route(clientsHandler, "GetClientInfo"))
	req, _ := http.NewRequest("GET", "/api/clients/"+client.Id+"/"+team.Alias, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"id":"backstage","secret":"SuperSecret","name":"Backstage","redirect_uri":"http://example.org/auth","owner":"owner@example.org","team":"team"}`)
}

func (s *S) TestGetClientInfoWhenClientIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Get("/api/clients/:id/:team", s.Api.Route(clientsHandler, "GetClientInfo"))
	req, _ := http.NewRequest("GET", "/api/clients/invalid-client/team", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Client not found on this team."}`)
}

func (s *S) TestGetClientInfoWhenIsNotInTeam(c *C) {
	bob.Save()
	owner.Save()
	team.Save(owner)
	client.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteClientByIdAndTeam(client.Id, team.Alias)
	defer bob.Delete()
	defer owner.Delete()

	s.router.Get("/api/clients/:id/:team", s.Api.Route(clientsHandler, "GetClientInfo"))
	req, _ := http.NewRequest("GET", "/api/clients/"+client.Id+"/"+team.Alias, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}
