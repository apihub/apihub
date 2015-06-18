package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/maestro/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestSaveWebhook(c *C) {
	team.Create(user)
	webookName := "Notify Slack."

	defer func() {
		webhook, _ := s.store.FindWebhookByName(webookName)
		webhook.Delete()
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(account.RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           "/api/webhooks",
		Body:           fmt.Sprintf(`{"name": "%s", "events": ["service.update"], "team": "%s", "config": {"url": "http://example.org"}}`, webookName, team.Alias),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"notify-slack","team":"backstage","events":["service.update"],"config":{"url":"http://example.org"}}`)
}

func (s *S) TestDeleteWebhookNotFound(c *C) {
	headers, code, body, err := httpClient.MakeRequest(account.RequestArgs{
		AcceptableCode: http.StatusNotFound,
		Method:         "DELETE",
		Path:           `/api/webhooks/not-found`,
		Body:           `{}`,
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Webhook not found."}`)
}
