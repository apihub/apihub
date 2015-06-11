package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	alias := "backstage-team"

	defer func() {
		store, _ := s.store()
		store.DeleteTeamByAlias(alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/teams",
		Body:    `{"name": "Backstage Team"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"Backstage Team","alias":"%s","users":["%s"],"owner":"%s"}`, alias, user.Email, user.Email))
}

func (s *S) TestCreateTeamWithCustomAlias(c *C) {
	alias := "backstage"

	defer func() {
		store, _ := s.store()
		store.DeleteTeamByAlias(alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/teams",
		Body:    fmt.Sprintf(`{"name": "Backstage Team", "alias": "%s"}`, alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"Backstage Team","alias":"%s","users":["%s"],"owner":"%s"}`, alias, user.Email, user.Email))
}

func (s *S) TestCreateTeamWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "POST", Path: "/api/teams", Body: `{"name": "Backstage Team"}`}, c)
}

func (s *S) TestCreateTeamWithInvalidRequest(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/teams",
		Body:    `"name": "Backstage Team"`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestTeamList(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    "/api/teams",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"items":[],"item_count":0}`)
}

func (s *S) TestTeamListWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "GET", Path: "/api/teams"}, c)
}

func (s *S) TestDeleteTeam(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/teams/%s", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"%s","alias":"%s","users":["%s"],"owner":"%s"}`, team.Name, team.Alias, user.Email, user.Email))
}

func (s *S) TestDeleteTeamWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/teams/backstage"}, c)
}

func (s *S) TestDeleteTeamWithoutPermission(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(alice)
	defer func() {
		store, _ := s.store()
		store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/teams/%s", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"Only the owner has permission to perform this operation."}`)
}

func (s *S) TestDeleteTeamNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/teams/not-found",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}
