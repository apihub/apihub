package api

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateService(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer account.DeleteServiceBySubdomain("backstage")

	payload := `{"subdomain": "backstage", "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams/:team/services", s.Api.route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/teams/"+team.Alias+"/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"subdomain":"backstage","description":"Useful desc.","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","transformers":[],"middlewares":[],"owner":"owner@example.org","team":"team","timeout":10}`
	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
	c.Assert(s.recorder.Body.String(), Equals, expected)
}

func (s *S) TestCreateServiceWhenUserIsNotSignedIn(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams/:team/services", s.Api.route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/teams/"+team.Alias+"/services", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestCreateServiceWhenTeamDoesNotExist(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	payload := `{"subdomain": "backstage", "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams/:team/services", s.Api.route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/teams/invalid-team/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Team not found."}`)
}

func (s *S) TestCreateServiceWithInvalidPayloadFormat(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	payload := `"subdomain": "backstage"`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams/:team/services", s.Api.route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/teams/invalid-team/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestDeleteService(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	s.router.Delete("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "DeleteService"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","description":"","disabled":true,"documentation":"","endpoint":"http://example.org/api","transformers":[],"middlewares":[],"owner":"owner@example.org","team":"team","timeout":0}`)
}

func (s *S) TestDeleteServiceWhenUserIsNotOwner(c *C) {
	alice.Save()
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer alice.Delete()
	defer owner.Delete()
	defer service.Delete()

	s.router.Delete("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "DeleteService"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestDeleteServiceIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Delete("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "DeleteService"))
	req, _ := http.NewRequest("DELETE", "/api/teams/"+team.Alias+"/services/invalid-service", nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestGetServiceInfo(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	s.router.Get("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "GetServiceInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","transformers":[],"middlewares":[],"owner":"owner@example.org","team":"team","timeout":0}`)
}

func (s *S) TestGetServiceInfoWhenServiceIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Get("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "GetServiceInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestGetServiceInfoWhenIsNotInTeam(c *C) {
	bob.Save()
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer bob.Delete()
	defer owner.Delete()
	defer service.Delete()

	s.router.Get("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "GetServiceInfo"))
	req, _ := http.NewRequest("GET", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}
