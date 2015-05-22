package gateway

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/gateway/middleware"
	"github.com/fatih/structs"
	. "gopkg.in/check.v1"
)

func (s *S) TestGatewayNotFound(c *C) {
	gateway := NewGateway(s.Settings)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "invalid.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusNotFound)
	c.Assert(w.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(w.Body.String(), Equals, "{\"error\":\"not_found\",\"error_description\":\"The requested resource could not be found but may be available again in the future. \"}\n")
}

func (s *S) TestGatewayExistingService(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		c.Assert(len(r.Header.Get("X-Request-Id")), Not(Equals), 0)
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}}
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
}

func (s *S) TestGatewayTimeout(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Timeout: 1}}
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusGatewayTimeout)
	c.Assert(w.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(w.Body.String(), Equals, "{\"error\":\"gateway_timeout\",\"error_description\":\"The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server.\"}")
}

func (s *S) TestGatewayCopyResponseHeaders(c *C) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test"}}
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusCreated)
	c.Assert(w.Header().Get("Content-Type"), Equals, "text/plain; charset=utf-8")
}

func (s *S) TestGatewayInternalError(c *C) {
	services := []*account.Service{&account.Service{Endpoint: "http://invalidurl", Subdomain: "test"}}
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
	c.Assert(w.Header().Get("Content-Type"), Equals, "application/json")
	c.Assert(w.Body.String(), Equals, "{\"error\":\"internal_server_error\",\"error_description\":\"dial tcp: lookup invalidurl: no such host\"}")
}

func (s *S) TestGatewayWithTransformer(c *C) {
	addBackstageHeader := func(r *http.Request, w *http.Response, buf *bytes.Buffer) {
		w.Header.Set("X-Backstage-Header", "Custom Header")
	}

	addViaHeader := func(r *http.Request, w *http.Response, buf *bytes.Buffer) {
		w.Header.Set("Via", "test.backstage.dev")
	}

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	services := []*account.Service{&account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Transformers: []string{"AddHeader", "AddHeaderVia"}}}
	gateway := NewGateway(s.Settings)
	gateway.Transformer().Add("AddHeader", addBackstageHeader)
	gateway.Transformer().Add("AddHeaderVia", addViaHeader)
	gateway.LoadServices(services)
	defer gateway.Close()

	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
	c.Assert(w.Header().Get("X-Backstage-Header"), Equals, "Custom Header")
	c.Assert(w.Header().Get("Via"), Equals, "test.backstage.dev")
}

func (s *S) TestAuthenticationMiddlewareWithoutHeader(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	owner.Save()
	defer owner.Delete()
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	service = &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Transformers: []string{"AddHeader", "AddHeaderVia"}}
	services := []*account.Service{service}

	service.Save(owner, team)
	defer service.Delete()

	conf := configAuthenticationMiddleware(service, owner)
	defer conf.Delete(owner)
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	gateway.ServeHTTP(w, r)

	c.Assert(w.Body.String(), Equals, `{"error":"unauthorized_access","error_description":"Request refused or access is not allowed."}`)
	c.Assert(w.Code, Equals, http.StatusUnauthorized)
}

func (s *S) TestAuthenticationMiddlewareWithInvalidHeader(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	service = &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Transformers: []string{"AddHeader", "AddHeaderVia"}, Team: team.Alias}
	service.Save(owner, team)
	defer service.Delete()

	services := []*account.Service{service}

	conf := configAuthenticationMiddleware(service, owner)
	defer conf.Delete(owner)
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	r.Header.Set("Authorization", "non-existing-token")
	gateway.ServeHTTP(w, r)

	c.Assert(w.Body.String(), Equals, `{"error":"unauthorized_access","error_description":"Request refused or access is not allowed."}`)
	c.Assert(w.Code, Equals, http.StatusUnauthorized)
}

func (s *S) TestAuthenticationMiddlewareWithValidHeader(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	auth := &api.AuthenticationInfo{ClientId: "123", Token: "test-123", Type: "Bearer", User: "alice@example.org", Expires: 10}
	s.AddToken(auth.Token, auth.Expires, structs.Map(auth))
	defer s.DeleteToken(auth.Token)

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Header.Get("Backstage-User"), Equals, auth.User)
		c.Assert(r.Header.Get("Backstage-ClientId"), Equals, auth.ClientId)
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	service = &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Transformers: []string{"AddHeader", "AddHeaderVia"}}
	services := []*account.Service{service}

	service.Save(owner, team)
	defer service.Delete()

	conf := configAuthenticationMiddleware(service, owner)
	defer conf.Delete(owner)
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	r.Header.Set("Authorization", auth.Token)
	gateway.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
}

func (s *S) TestAuthenticationMiddlewareWithValidHeaderForApp(c *C) {
	owner.Save()
	team.Save(owner)
	defer owner.Delete()
	defer account.DeleteTeamByAlias(team.Alias, owner)

	auth := &api.AuthenticationInfo{ClientId: "123", Token: "test-123", Type: "Bearer", Expires: 10}
	s.AddToken(auth.Token, auth.Expires, structs.Map(auth))
	defer s.DeleteToken(auth.Token)

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.Header["Backstage-User"], IsNil)
		c.Assert(r.Header.Get("Backstage-ClientId"), Equals, auth.ClientId)
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	service = &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Transformers: []string{"AddHeader", "AddHeaderVia"}}
	service.Save(owner, team)
	defer service.Delete()
	services := []*account.Service{service}

	conf := configAuthenticationMiddleware(service, owner)
	defer conf.Delete(owner)
	gateway := NewGateway(s.Settings)
	gateway.LoadServices(services)
	defer gateway.Close()
	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://test.backstage.dev", nil)
	r.Header.Set("Authorization", auth.Token)
	gateway.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
}

func (s *S) TestHasMiddleware(c *C) {
	gateway := NewGateway(s.Settings)
	defer gateway.Close()
	c.Assert(gateway.HasMiddleware("oauth"), Equals, false)

	gateway.Middleware().Add("oauth", middleware.NewAuthenticationMiddleware)
	c.Assert(gateway.HasMiddleware("oauth"), Equals, true)
}

func configAuthenticationMiddleware(s *account.Service, user *account.User) *account.PluginConfig {
	conf := &account.PluginConfig{
		Name:    "authentication",
		Service: s.Subdomain,
		Config:  map[string]interface{}{},
	}

	err := conf.Save(user)
	if err != nil {
		panic(err)
	}
	return conf
}
