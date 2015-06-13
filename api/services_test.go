package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/backstage/account"
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

func (s *S) TestUpdateService(c *C) {
	team.Create(user)
	service.Create(user, team)

	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/services/%s", service.Subdomain),
		Body:    fmt.Sprintf(`{"documentation": "http://docs.org", "disabled": true}`, team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"subdomain":"backstage","disabled":true,"documentation":"http://docs.org","endpoint":"http://example.org/api","owner":"bob@bar.example.org","team":"backstage","timeout":0}`)
}

func (s *S) TestUpdateServiceNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    "/api/services/not_found",
		Body:    `{}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestUpdateServiceNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	t := account.Team{Name: "example"}
	t.Create(alice)
	service.Create(alice, t)
	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(t.Alias)
		alice.Delete()
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/services/%s", service.Subdomain),
		Body:    `{}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestUpdateServiceWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "PUT", Path: "/api/services/subdomain", Body: `{}`}, c)
}

func (s *S) TestDeleteService(c *C) {
	service.Create(user, team)

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/services/%s", service.Subdomain),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"subdomain":"backstage","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"bob@bar.example.org","team":"backstage","timeout":0}`)
}

func (s *S) TestDeleteServiceWithoutPermission(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	service.Create(alice, team)
	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/services/%s", service.Subdomain),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"Only the owner has permission to perform this operation."}`)
}

func (s *S) TestDeleteServiceIsNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/services/not-found",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestDeleteServiceWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/services/subdomain", Body: `{}`}, c)
}

func (s *S) TestServiceInfo(c *C) {
	team.Create(user)
	service.Create(user, team)
	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/services/%s", service.Subdomain),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"subdomain":"backstage","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"bob@bar.example.org","team":"backstage","timeout":0}`)
}

func (s *S) TestServiceInfoNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	t := account.Team{Name: "example"}
	t.Create(alice)
	service.Create(alice, t)
	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(t.Alias)
		alice.Delete()
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/services/%s", service.Subdomain),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestServiceInfoNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    "/api/services/not-found",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestServiceList(c *C) {
	team.Create(user)
	service.Create(user, team)
	defer func() {
		store, _ := s.store()
		serv, _ := store.FindServiceBySubdomain(service.Subdomain)
		store.DeleteService(serv)
		store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    "/api/services",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"items":[{"subdomain":"backstage","disabled":false,"documentation":"","endpoint":"http://example.org/api","owner":"bob@bar.example.org","team":"backstage","timeout":0}],"item_count":1}`)
}

func (s *S) TestServiceListWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "GET", Path: "/api/services", Body: `{}`}, c)
}
