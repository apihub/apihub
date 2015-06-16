package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/apimanager/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	alias := "backstage-team"

	defer func() {
		s.store.DeleteTeamByAlias(alias)
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
		s.store.DeleteTeamByAlias(alias)
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

func (s *S) TestCreateTeamWhenAlreadyExists(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)

	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/teams",
		Body:    fmt.Sprintf(`{"name": "Backstage Team", "alias": "%s"}`, team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Someone already has that team alias. Could you try another?"}`)

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
		s.store.DeleteTeamByAlias(team.Alias)
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

func (s *S) TestTeamInfo(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)
	defer team.Delete(user)

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/teams/%s", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"%s","alias":"%s","users":["%s"],"owner":"%s"}`, team.Name, team.Alias, user.Email, user.Email))
	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
}

func (s *S) TestTeamInfoNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    "/api/teams/not-found",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}
func (s *S) TestTeamInfoWithoutPermission(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(alice)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/teams/%s", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestAddUser(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/teams/%s/users", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
		Body:    fmt.Sprintf(`{"users": ["%s"]}`, alice.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"Backstage Team","alias":"backstage","users":["bob@bar.example.org","alice@bar.example.org"],"owner":"bob@bar.example.org"}`)
}

func (s *S) TestAddUserNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(alice)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/teams/%s/users", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
		Body:    fmt.Sprintf(`{"users": ["%s"]}`, alice.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestAddUserWithoutSignIn(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	testWithoutSignIn(RequestArgs{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/teams/%s/users", team.Alias),
		Body:   `{"users": ["bob@example.org"]}`},
		c)
}

func (s *S) TestAddUserNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    "/api/teams/not-found/users",
		Body:    `{"name": "New name"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestRemoveUser(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	team := account.Team{Name: "Backstage Team", Alias: "backstage", Users: []string{alice.Email}}
	team.Create(user)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/teams/%s/users", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
		Body:    fmt.Sprintf(`{"users": ["%s"]}`, alice.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"Backstage Team","alias":"backstage","users":["bob@bar.example.org"],"owner":"bob@bar.example.org"}`)
}

func (s *S) TestRemoveUserWithoutSignIn(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	testWithoutSignIn(RequestArgs{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/teams/%s/users", team.Alias),
		Body:   `{"users": ["bob@example.org"]}`},
		c)
}

func (s *S) TestRemoveUserNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(alice)
	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/teams/%s/users", team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
		Body:    fmt.Sprintf(`{"users": ["%s"]}`, alice.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestRemoveUserNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/teams/not-found/users",
		Body:    `{"name": "New name"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestUpdateTeam(c *C) {
	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(user)

	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/teams/%s", team.Alias),
		Body:    `{"name": "New name"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"New name","alias":"%s","users":["%s"],"owner":"%s"}`, team.Alias, user.Email, user.Email))
}

func (s *S) TestUpdateTeamNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	team := account.Team{Name: "Backstage Team", Alias: "backstage"}
	team.Create(alice)

	defer func() {
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/teams/%s", team.Alias),
		Body:    `{"name": "New name"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
}

func (s *S) TestUpdateTeamNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    "/api/teams/not-found",
		Body:    `{"name": "New name"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}
