package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestAddGetRequestError(c *C) {
	m := web.New()

	m.Get("/helloworld", func(c web.C, w http.ResponseWriter, r *http.Request) {
		AddRequestError(&c, &errors.HTTPError{StatusCode: http.StatusUnauthorized,
			Message: "You do not have access to this resource."})

		key, _ := GetRequestError(&c)
		body, _ := json.Marshal(key)
		http.Error(w, string(body), key.StatusCode)
	})

	req, _ := http.NewRequest("GET", "/helloworld", nil)
	recorder := httptest.NewRecorder()
	env := map[string]interface{}{}
	m.ServeHTTPC(web.C{Env: env}, recorder, req)

	c.Assert(recorder.Code, Equals, 401)
	c.Assert(recorder.Body.String(), Equals, "{\"status_code\":401,\"message\":\"You do not have access to this resource.\",\"url\":\"\"}\n")
}

func (s *S) TestSetAndGetCurrentUser(c *C) {
	m := web.New()

	m.Get("/helloworld", func(c web.C, w http.ResponseWriter, r *http.Request) {
		alice := &account.User{Username: "alice", Name: "Alice", Email: "alice@example.org", Password: "123456"}
		alice.Save()
		defer alice.Delete()
		SetCurrentUser(&c, alice)
		user, _ := GetCurrentUser(&c)
		body, _ := json.Marshal(user)
		fmt.Fprint(w, string(body))
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/helloworld", nil)
	recorder := httptest.NewRecorder()
	env := map[string]interface{}{}
	m.ServeHTTPC(web.C{Env: env}, recorder, req)

	c.Assert(recorder.Code, Equals, 200)
}

func (s *S) TestGetCurrentUserWhenNotSignedIn(c *C) {
	m := web.New()

	m.Get("/helloworld", func(co web.C, w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentUser(&co)
		c.Assert(err.Error(), Equals, "User is not signed in.")
	})

	req, _ := http.NewRequest("GET", "/helloworld", nil)
	recorder := httptest.NewRecorder()
	env := map[string]interface{}{}
	m.ServeHTTPC(web.C{Env: env}, recorder, req)
}
