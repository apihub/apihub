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

	payload := `{"subdomain": "backstage", "team": "` + team.Alias + `", "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10, "transformers": ["test"]}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("POST", "/api/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"subdomain":"backstage","description":"Useful desc.","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","transformers":["test"],"owner":"owner@example.org","team":"team","timeout":10}`
	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
	c.Assert(s.recorder.Body.String(), Equals, expected)
}

func (s *S) TestCreateServiceWhenAlreadyExists(c *C) {
	alice.Save()
	owner.Save()
	team.Save(owner)
	service.Save(owner, team)
	aliceTeam := &account.Team{Name: "Alice Team"}
	aliceTeam.Save(alice)
	defer account.DeleteTeamByAlias(team.Alias, owner)
	defer account.DeleteTeamByAlias(aliceTeam.Alias, alice)
	defer owner.Delete()
	defer alice.Delete()
	defer account.DeleteServiceBySubdomain("backstage")

	payload := `{"subdomain": "backstage", "team": "` + aliceTeam.Alias + `", "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10, "transformers": ["test"]}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("POST", "/api/services", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"There is another service with this subdomain."}`)
}

func (s *S) TestCreateServiceWhenUserIsNotSignedIn(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("POST", "/api/services", b)
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

	payload := `{"subdomain": "backstage", "team": "invalid-team", "description": "Useful desc.", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("POST", "/api/services", b)
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

	req, _ := http.NewRequest("POST", "/api/services", b)
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

	payload := `{"description": "New DESC", "team":"` + team.Alias + `", "disabled": true, "documentation": "http://backstage.example.org/doc", "endpoint": "http://github.com/backstage/backstage", "timeout": 1}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("PUT", "/api/services/"+service.Subdomain, b)
	req.Header.Set("Content-Type", "application/json")
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	expected := `{"subdomain":"backstage","description":"New DESC","disabled":true,"documentation":"http://backstage.example.org/doc","endpoint":"http://github.com/backstage/backstage","owner":"owner@example.org","team":"` + team.Alias + `","timeout":1}`
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

	req, _ := http.NewRequest("PUT", "/api/services/"+service.Subdomain, b)
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

	payload := `{"team": "invalid-team"}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("PUT", "/api/services/"+service.Subdomain, b)
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

	req, _ := http.NewRequest("PUT", "/api/services/"+service.Subdomain, b)
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

	payload := `{"team":"` + team.Alias + `"}`
	b := strings.NewReader(payload)

	req, _ := http.NewRequest("PUT", "/api/services/"+service.Subdomain, b)
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

	req, _ := http.NewRequest("DELETE", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
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

	req, _ := http.NewRequest("DELETE", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestDeleteServiceIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	req, _ := http.NewRequest("DELETE", "/api/services/invalid-service", nil)
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

	req, _ := http.NewRequest("GET", "/api/services/"+service.Subdomain, nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"subdomain":"backstage","description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}`)
}

func (s *S) TestGetServiceInfoWhenServiceIsNotFound(c *C) {
	bob.Save()
	defer bob.Delete()

	req, _ := http.NewRequest("GET", "/api/services/"+service.Subdomain, nil)
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

	req, _ := http.NewRequest("GET", "/api/services/"+service.Subdomain, nil)
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

	req, _ := http.NewRequest("GET", "/api/services", nil)
	s.env[CurrentUser] = owner
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"items":[{"subdomain":"_get_user_services","description":"","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"owner@example.org","team":"team","timeout":0}],"item_count":1}`)
}

func (s *S) TestGetUserServicesWhenUserIsNotSignedIn(c *C) {
	req, _ := http.NewRequest("GET", "/api/services", nil)
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}
