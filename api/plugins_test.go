package api

import (
  "net/http"
  "strings"

  "github.com/backstage/backstage/account"
  "github.com/zenazn/goji/web"
  . "gopkg.in/check.v1"
)


func (s *S) TestSubscribePlugin(c *C) {
  owner.Save()
  team.Save(owner)
  service.Save(owner, team)
  defer owner.Delete()
  defer account.DeleteTeamByAlias(team.Alias, owner)
  defer account.DeleteServiceBySubdomain(service.Subdomain)

  payload := `{"service":"` + service.Subdomain + `","config":{"timeout":123}}`
  b := strings.NewReader(payload)

  s.router.Put("/api/plugins/:name/subscriptions", s.Api.route(pluginsHandler, "SubscribePlugin"))
  req, _ := http.NewRequest("PUT", "/api/plugins/cors/subscriptions", b)
  req.Header.Set("Content-Type", "application/json")
  s.env[CurrentUser] = owner
  webC := web.C{Env: s.env}
  s.router.ServeHTTPC(webC, s.recorder, req)

  expected := `{"name":"cors","service":"backstage","config":{"timeout":123}}`
  c.Assert(s.recorder.Header().Get("Content-Type"), Equals, "application/json")
  c.Assert(s.recorder.Body.String(), Equals, expected)
  c.Assert(s.recorder.Code, Equals, http.StatusOK)
}

func (s *S) TestSubscribePluginWhenUserIsNotSignedIn(c *C) {
  payload := `{}`
  b := strings.NewReader(payload)

  s.router.Put("/api/plugins/:name/subscriptions", s.Api.route(pluginsHandler, "SubscribePlugin"))
  req, _ := http.NewRequest("PUT", "/api/plugins/cors/subscriptions", b)
  req.Header.Set("Content-Type", "application/json")
  webC := web.C{Env: s.env}
  s.router.ServeHTTPC(webC, s.recorder, req)

  c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
  c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestSubscribePluginWhenServiceDoesNotExist(c *C) {
  owner.Save()
  team.Save(owner)
  defer account.DeleteTeamByAlias(team.Alias, owner)
  defer owner.Delete()

  payload := `{"service":"invalidservice","config":{"timeout":123}}`
  b := strings.NewReader(payload)

  s.router.Put("/api/plugins/:name/subscriptions", s.Api.route(pluginsHandler, "SubscribePlugin"))
  req, _ := http.NewRequest("PUT", "/api/plugins/cors/subscriptions", b)
  req.Header.Set("Content-Type", "application/json")
  s.env[CurrentUser] = owner
  webC := web.C{Env: s.env}
  s.router.ServeHTTPC(webC, s.recorder, req)

  c.Assert(s.recorder.Code, Equals, http.StatusNotFound)
  c.Assert(s.recorder.Body.String(), Equals, `{"error":"not_found","error_description":"Service not found."}`)
}
