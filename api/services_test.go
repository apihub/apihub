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

	payload := `{"subdomain": "backstage", "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10, "transformers": ["test"]}`
	b := strings.NewReader(payload)

	s.router.Post("/api/teams/:team/services", s.Api.route(servicesHandler, "CreateService"))
	req, _ := http.NewRequest("POST", "/api/teams/"+team.Alias+"/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"subdomain":"backstage","description":"Useful desc.","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","transformers":["test"],"owner":"owner@example.org","team":"team","timeout":10}`
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

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
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

func (s *S) TestUpdateService(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	payload := `{"subdomain": "backstage", "description": "New DESC", "disabled": true, "documentation": "http://backstage.example.org/doc", "endpoint": "http://github.com/backstage/backstage", "timeout": 1}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "UpdateService"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"subdomain":"backstage","description":"New DESC","disabled":true,"documentation":"http://backstage.example.org/doc","endpoint":"http://github.com/backstage/backstage","owner":"owner@example.org","team":"team","timeout":1}`
	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, expected)
}

func (s *S) TestUpdateServiceWhenUserIsNotSignedIn(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "UpdateService"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestUpdateServiceWhenTeamDoesNotExist(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "UpdateService"))
	req, _ := http.NewRequest("PUT", "/api/teams/invalid-team/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestUpdateServiceWithInvalidPayloadFormat(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()
	defer service.Delete()

	payload := `"subdomain": "backstage"`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "UpdateService"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestUpdateServiceWhenDoesNotBelongToTeam(c *C) {
	alice.Save()
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer alice.Delete()
	defer owner.Delete()
	defer service.Delete()

	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "UpdateService"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusForbidden)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
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
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","description":"","disabled":true,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}`)
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
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}`)
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

func (s *S) TestGetUserServices(c *C) {
	owner.Save()
	team.Save(owner)
	service := &account.Service{Endpoint: "http://example.org/api", Subdomain: "_get_user_services", Transformers: []string{}}
	service.Save(owner, team)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteServiceBySubdomain(service.Subdomain)

	s.router.Get("/api/teams/services", s.Api.route(servicesHandler, "GetUserServices"))
	req, _ := http.NewRequest("GET", "/api/teams/services", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"items":[{"subdomain":"_get_user_services","description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}],"item_count":1}`)
}

func (s *S) TestGetUserServicesWhenUserIsNotSignedIn(c *C) {
	s.router.Get("/api/teams/services", s.Api.route(servicesHandler, "GetUserServices"))
	req, _ := http.NewRequest("GET", "/api/teams/services", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestConfigurePlugin(c *C) {
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteServiceBySubdomain(service.Subdomain)

	payload := `{"name":"cors","config":{"timeout":123}}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "ConfigurePlugin"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"name":"cors","config":{"timeout":123}}`
	c.Assert(s.recorder.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(s.recorder.Body.String(), Equals, expected)
	c.Assert(s.recorder.Code, Equals, http.StatusOK)
}

func (s *S) TestConfigurePluginWhenUserIsNotSignedIn(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "ConfigurePlugin"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestConfigurePluginWhenTeamDoesNotExist(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	payload := `{"name":"cors","config":{"timeout":123}}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "ConfigurePlugin"))
	req, _ := http.NewRequest("PUT", "/api/teams/invalid-team/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestConfigurePluginWhenServiceDoesNotExist(c *C) {
	owner.Save()
	team.Save(owner)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer owner.Delete()

	payload := `{"name":"cors","config":{"timeout":123}}`
	b := strings.NewReader(payload)

	s.router.Put("/api/teams/:team/services/:subdomain", s.Api.route(servicesHandler, "ConfigurePlugin"))
	req, _ := http.NewRequest("PUT", "/api/teams/"+team.Alias+"/services/invalid-service", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}
