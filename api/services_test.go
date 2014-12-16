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

	payload := `{"team": "` + team.Alias + `", "subdomain": "backstage", "allow_keyless_use": true, "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}`
	b := strings.NewReader(payload)

	s.router.Post("/api/services", s.Api.Route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"subdomain":"backstage","allow_keyless_use":true,"description":"Useful desc.","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","owner":"owner@example.org","team":"team","timeout":10}`
	c.Assert(s.recorder.Code, Equals, 201)
	c.Assert(s.recorder.Body.String(), Equals, expected)
}

func (s *S) TestCreateServiceWhenUserIsNotSignedIn(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Post("/api/services", s.Api.Route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/services", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestCreateServiceWhenTeamDoesNotExist(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer account.DeleteServiceBySubdomain("backstage")

	payload := `{"team": "invalid-team", "subdomain": "backstage", "allow_keyless_use": true, "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}`
	b := strings.NewReader(payload)

	s.router.Post("/api/services", s.Api.Route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"Team not found."}`)
}

func (s *S) TestCreateServiceWithInvalidMessageFormat(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer account.DeleteServiceBySubdomain("backstage")

	payload := `"subdomain": "backstage"`
	b := strings.NewReader(payload)

	s.router.Post("/api/services", s.Api.Route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"The request was invalid or cannot be served."}`)
}

func (s *S) TestDeleteService(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	s.router.Delete("/api/services/:subdomain", s.Api.Route(servicesHandler, "DeleteService"))
	req, _ := http.NewRequest("DELETE", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","allow_keyless_use":false,"description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}`)
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

	s.router.Delete("/api/teams/:alias", s.Api.Route(teamsHandler, "DeleteTeam"))
	s.router.Delete("/api/services/:subdomain", s.Api.Route(servicesHandler, "DeleteService"))
	req, _ := http.NewRequest("DELETE", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"Service not found or you're not the owner."}`)
}

func (s *S) TestDeleteServiceIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Delete("/api/teams/:alias", s.Api.Route(teamsHandler, "DeleteTeam"))
	s.router.Delete("/api/services/:subdomain", s.Api.Route(servicesHandler, "DeleteService"))
	req, _ := http.NewRequest("DELETE", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"Service not found or you're not the owner."}`)
}

func (s *S) TestGetServiceInfo(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	s.router.Get("/api/services/:subdomain", s.Api.Route(servicesHandler, "GetServiceInfo"))
	req, _ := http.NewRequest("GET", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","allow_keyless_use":false,"description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}`)
}

func (s *S) TestGetServiceInfoWhenServiceIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	s.router.Get("/api/services/:subdomain", s.Api.Route(servicesHandler, "GetServiceInfo"))
	req, _ := http.NewRequest("GET", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"Service not found or you dont belong to the team responsible for it."}`)
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

	s.router.Get("/api/services/:subdomain", s.Api.Route(servicesHandler, "GetServiceInfo"))
	req, _ := http.NewRequest("GET", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = bob
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 403)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":403,"message":"You do not belong to this team!"}`)
}
