package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/apimanager/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateApp(c *C) {
	team.Create(user)
	clientId := "backstage"

	defer func() {
		app, _ := s.store.FindAppByClientId(clientId)
		s.store.DeleteApp(app)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/apps",
		Body:    fmt.Sprintf(`{"name": "Ios App", "client_id": "%s", "client_secret": "secret","redirect_uris": ["http://www.example.org/auth"], "team": "%s" }`, clientId, team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"client_id":"backstage","client_secret":"secret","name":"Ios App","redirect_uris":["http://www.example.org/auth"],"owner":"bob@bar.example.org","team":"backstage"}`)
}

func (s *S) TestCreateAppWhenAlreadyExists(c *C) {
	team.Create(user)
	app.Team = team.Alias
	app.Create(user, team)

	defer func() {
		app, _ := s.store.FindAppByClientId(app.ClientId)
		s.store.DeleteApp(app)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/apps",
		Body:    fmt.Sprintf(`{"name": "Ios App", "client_id": "%s", "client_secret": "secret","redirect_uris": ["http://www.example.org/auth"], "team": "%s" }`, app.ClientId, team.Alias),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"There is another app with this client id."}`)
}

func (s *S) TestCreateAppWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "POST", Path: "/api/apps", Body: `{}`}, c)
}

func (s *S) TestCreateAppTeamNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/apps",
		Body:    fmt.Sprintf(`{"name": "Ios App", "client_id": "%s", "client_secret": "secret","redirect_uris": ["http://www.example.org/auth"], "team": "not_found" }`, app.ClientId),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Team not found."}`)
}

func (s *S) TestCreateAppInvalidBody(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "POST",
		Path:    "/api/apps",
		Body:    "invalid:body",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestUpdateApp(c *C) {
	team.Create(user)
	app.Create(user, team)

	defer func() {
		ap, _ := s.store.FindAppByClientId(app.ClientId)
		s.store.DeleteApp(ap)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/apps/%s", app.ClientId),
		Body:    `{"name": "new name", "client_secret": "new secret"}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"client_id":"ios","client_secret":"new secret","name":"new name","redirect_uris":["http://www.example.org/auth"],"owner":"bob@bar.example.org","team":"backstage"}`)
}

func (s *S) TestUpdateAppNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    "/api/apps/not_found",
		Body:    `{}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"App not found."}`)
}

func (s *S) TestUpdateAppNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	t := account.Team{Name: "example"}
	t.Create(alice)
	app.Create(alice, t)
	defer func() {
		ap, _ := s.store.FindAppByClientId(app.ClientId)
		s.store.DeleteApp(ap)
		s.store.DeleteTeamByAlias(t.Alias)
		alice.Delete()
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "PUT",
		Path:    fmt.Sprintf("/api/apps/%s", app.ClientId),
		Body:    `{}`,
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestUpdateAppWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "PUT", Path: "/api/apps/client_id", Body: `{}`}, c)
}

func (s *S) TestDeleteApp(c *C) {
	app.Create(user, team)

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/apps/%s", app.ClientId),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"client_id":"ios","client_secret":"secret","name":"Ios App","redirect_uris":["http://www.example.org/auth"],"owner":"bob@bar.example.org","team":"backstage"}`)
}

func (s *S) TestDeleteAppWithoutPermission(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	defer alice.Delete()

	app.Create(alice, team)
	defer func() {
		ap, _ := s.store.FindAppByClientId(app.ClientId)
		s.store.DeleteApp(ap)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    fmt.Sprintf("/api/apps/%s", app.ClientId),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"Only the owner has permission to perform this operation."}`)
}

func (s *S) TestDeleteAppIsNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/apps/not-found",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"App not found."}`)
}

func (s *S) TestDeleteAppWithoutSignIn(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/apps/client_id", Body: `{}`}, c)
}

func (s *S) TestAppInfo(c *C) {
	team.Create(user)
	app.Create(user, team)

	defer func() {
		ap, _ := s.store.FindAppByClientId(app.ClientId)
		s.store.DeleteApp(ap)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/apps/%s", app.ClientId),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"client_id":"ios","client_secret":"secret","name":"Ios App","redirect_uris":["http://www.example.org/auth"],"owner":"bob@bar.example.org","team":"backstage"}`)
}

func (s *S) TestAppInfoNotMember(c *C) {
	alice := account.User{Name: "alice", Email: "alice@bar.example.org", Password: "secret"}
	alice.Create()
	t := account.Team{Name: "example"}
	t.Create(alice)
	app.Create(alice, t)

	defer func() {
		ap, _ := s.store.FindAppByClientId(app.ClientId)
		s.store.DeleteApp(ap)
		s.store.DeleteTeamByAlias(t.Alias)
		alice.Delete()
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/apps/%s", app.ClientId),
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusForbidden)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"access_denied","error_description":"You do not belong to this team!"}`)
}

func (s *S) TestAppInfoNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "GET",
		Path:    "/api/apps/not-found",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"App not found."}`)
}
