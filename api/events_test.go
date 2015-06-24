package api_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/requests"
	. "gopkg.in/check.v1"
)

func (s *S) TestListenEvents(c *C) {
	var wg sync.WaitGroup
	wg.Add(3)
	s.api.ListenEvents()

	// Custom Text
	customText := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
		body, _ := ioutil.ReadAll(r.Body)
		c.Assert(string(body), Equals, `{"username": "Backstage Maestro", "channel": "#backstage",
		"icon_url": "http://www.albertoleal.me/images/maestro-pq.png",
		"text": "Um novo serviço foi criado no Backstage Maestro, com o seguinte subdomínio: backstage."}`)
	})
	svrC := httptest.NewServer(customText)
	s.api.AddHook(account.Hook{
		Name:   "backstage-maestro-custom",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: svrC.URL},
		Text: `{"username": "Backstage Maestro", "channel": "#backstage",
		"icon_url": "http://www.albertoleal.me/images/maestro-pq.png",
		"text": "Um novo serviço foi criado no Backstage Maestro, com o seguinte subdomínio: {{.Service.Subdomain}}."}`,
	})

	// Default Text
	defaultText := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
		body, _ := ioutil.ReadAll(r.Body)
		c.Assert(string(body), Matches, fmt.Sprintf(`{"created_at":".*","name":"service.create","service":{"subdomain":"backstage","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"backstage"}}`))
	})
	svrD := httptest.NewServer(defaultText)
	s.api.AddHook(account.Hook{
		Name:   "backstage-maestro-default",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: svrD.URL},
	})

	// With wrong Template
	wrongTmpl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
		body, _ := ioutil.ReadAll(r.Body)
		c.Assert(string(body), Matches, fmt.Sprintf(`{"created_at":".*","name":"service.create","service":{"subdomain":"backstage","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"backstage"}}`))
	})
	svrT := httptest.NewServer(wrongTmpl)
	s.api.AddHook(account.Hook{
		Name:   "backstage-maestro-wrong-tmpl",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: svrT.URL},
		Text:   `{"username": "Backstage Maestro", "channel": "{{.Team.NotFound}}"}`,
	})

	team.Create(user)
	subdomain := "backstage"

	defer func() {
		serv, _ := s.store.FindServiceBySubdomain(subdomain)
		s.store.DeleteService(serv)
		s.store.DeleteTeamByAlias(team.Alias)
	}()

	headers, code, body, _ := httpClient.MakeRequest(requests.Args{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/services",
		Body:           fmt.Sprintf(`{"subdomain": "%s", "endpoint": "http://example.org", "team": "%s"}`, subdomain, team.Alias),
		Headers:        http.Header{"Authorization": {s.authHeader}},
	})

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"subdomain":"backstage","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"backstage"}`)
	wg.Wait()
}
