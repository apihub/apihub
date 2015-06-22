package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/maestro/requests"
	. "gopkg.in/check.v1"
)

func (s *S) TestSaveHook(c *C) {
	team.Create(user)
	webookName := "Notify Slack."

	defer func() {
		hook, _ := s.store.FindHookByName(webookName)
		hook.Delete()
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           "/api/hooks",
		Body:           fmt.Sprintf(`{"name": "%s", "events": ["service.update"], "team": "%s", "config": {"address": "http://example.org"}}`, webookName, team.Alias),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"notify-slack","team":"backstage","events":["service.update"],"config":{"address":"http://example.org"}}`)
}

func (s *S) TestDeleteHookNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusNotFound,
		Method:         "DELETE",
		Path:           `/api/hooks/not-found`,
		Body:           `{}`,
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Hook not found."}`)
}
