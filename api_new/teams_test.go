package api_new_test

import (
	"fmt"
	"net/http"

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
