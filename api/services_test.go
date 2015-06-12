package api_test

import (
	"fmt"
	"net/http"

	// "github.com/backstage/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateService(c *C) {
	team.Create(user)
	subdomain := "backstage"

	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/services",
		Body:    fmt.Sprintf(`{"subdomain": "%s", "endpoint": "http://example.org", "team": "%s"}`, subdomain, team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"subdomain":"backstage","disabled":false,"documentation":"","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"backstage","timeout":0}`)
}

func (s *S) TestCreateServiceWhenAlreadyExists(c *C) {
	team.Create(user)
	service.Team = team.Alias
	service.Create(user, team)

	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/services",
		Body:    fmt.Sprintf(`{"subdomain": "%s", "endpoint": "http://example.org", "team": "%s"}`, service.Subdomain, team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"There is another service with this subdomain."}`)
}

func (s *S) TestCreateServiceWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "POST", Path: "/api/services", Body: `{}`}, c)
}

func (s *S) TestCreateServiceTeamNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/services",
		Body:    fmt.Sprintf(`{"subdomain": "%s", "endpoint": "http://example.org", "team": "not_found"}`, service.Subdomain),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestCreateServiceInvalidBody(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/services",
		Body:    "invalid:body",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}
