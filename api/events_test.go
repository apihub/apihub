package api_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/requests"
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
		c.Assert(string(body), Equals, `{"username": "ApiHub", "channel": "#apihub",
		"icon_url": "http://www.albertoleal.me/images/apihub-pq.png",
		"text": "Um novo serviço foi criado no ApiHub, com o seguinte subdomínio: apihub."}`)
	})
	svrC := httptest.NewServer(customText)
	s.api.AddHook(account.Hook{
		Name:   "apihub-apihub-custom",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: svrC.URL},
		Text: `{"username": "ApiHub", "channel": "#apihub",
		"icon_url": "http://www.albertoleal.me/images/apihub-pq.png",
		"text": "Um novo serviço foi criado no ApiHub, com o seguinte subdomínio: {{.Service.Subdomain}}."}`,
	})

	// Default Text
	defaultText := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
		body, _ := ioutil.ReadAll(r.Body)
		c.Assert(string(body), Matches, fmt.Sprintf(`{"created_at":".*","name":"service.create","service":{"subdomain":"apihub","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"apihub"}}`))
	})
	svrD := httptest.NewServer(defaultText)
	s.api.AddHook(account.Hook{
		Name:   "apihub-apihub-default",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: svrD.URL},
	})

	// With wrong Template
	wrongTmpl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
		body, _ := ioutil.ReadAll(r.Body)
		c.Assert(string(body), Matches, fmt.Sprintf(`{"created_at":".*","name":"service.create","service":{"subdomain":"apihub","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"apihub"}}`))
	})
	svrT := httptest.NewServer(wrongTmpl)
	s.api.AddHook(account.Hook{
		Name:   "apihub-apihub-wrong-tmpl",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: svrT.URL},
		Text:   `{"username": "ApiHub", "channel": "{{.Team.NotFound}}"}`,
	})

	team.Create(user)
	subdomain := "apihub"

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
	c.Assert(string(body), Equals, `{"subdomain":"apihub","endpoint":"http://example.org","owner":"bob@bar.example.org","team":"apihub"}`)
	wg.Wait()
}
