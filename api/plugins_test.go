package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/maestro/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestSubscribePlugin(c *C) {
	team.Create(user)
	service.Create(user, team)
	pluginName := "cors"

	defer func() {
		serv, _ := s.store.FindServiceBySubdomain(service.Subdomain)
		s.store.DeleteService(serv)
		s.store.DeleteTeamByAlias(team.Alias)

		plugin, _ := s.store.FindPluginConfigByNameAndService(pluginName, service)
		s.store.DeletePluginConfig(plugin)
	}()

	headers, code, body, err := httpClient.MakeRequest(account.RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           fmt.Sprintf(`/api/services/%s/plugins`, service.Subdomain),
		Body:           fmt.Sprintf(`{"name": "%s", "config": {"version": 1}}`, pluginName),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"cors","service":"backstage","config":{"version":1}}`)
}

func (s *S) TestSubscribePluginNotFound(c *C) {
	pluginName := "cors"

	headers, code, body, err := httpClient.MakeRequest(account.RequestArgs{
		AcceptableCode: http.StatusNotFound,
		Method:         "PUT",
		Path:           `/api/services/not-found/plugins`,
		Body:           fmt.Sprintf(`{"name": "%s", "config": {"version": 1}}`, pluginName),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}

func (s *S) TestUnsubscribePlugin(c *C) {
	team.Create(user)
	service.Create(user, team)
	pluginConfig.Save(service)

	defer func() {
		serv, _ := s.store.FindServiceBySubdomain(service.Subdomain)
		s.store.DeleteService(serv)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(account.RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf(`/api/services/%s/plugins/%s`, service.Subdomain, pluginConfig.Name),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"Plugin Config","service":"backstage","config":{"version":1}}`)
}

func (s *S) TestUnsubscribePluginNotFound(c *C) {
	team.Create(user)
	service.Create(user, team)

	defer func() {
		serv, _ := s.store.FindServiceBySubdomain(service.Subdomain)
		s.store.DeleteService(serv)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, err := httpClient.MakeRequest(account.RequestArgs{
		AcceptableCode: http.StatusNotFound,
		Method:         "DELETE",
		Path:           fmt.Sprintf(`/api/services/%s/plugins/not-found`, service.Subdomain),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"not_found","error_description":"Plugin Config not found."}`)
}
