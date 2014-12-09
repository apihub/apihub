package api

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	defer func() {
		user, err := account.FindUserByEmail("alice@example.org")
		c.Assert(err, IsNil)
		err = user.Delete()
		c.Assert(err, IsNil)
	}()
	payload := `{"name": "Alice", "email": "alice@example.org", "username": "alice", "password": "123456"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/users", s.Api.Route(usersHandler, "CreateUser"))
	req, _ := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 201)
  c.Assert(s.recorder.Body.String(), Equals, `{"name":"Alice","email":"alice@example.org","username":"alice"}`)
}

func (s *S) TestCreateUserWithInvalidPayloadFormat(c *C) {
	payload := `"name": "Alice"`
	b := strings.NewReader(payload)

	s.router.Post("/api/users", s.Api.Route(usersHandler, "CreateUser"))
	req, _ := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"The request was bad-formed."}`)
}

func (s *S) TestCreateUserWithMissingRequiredFields(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Post("/api/users", s.Api.Route(usersHandler, "CreateUser"))
	req, _ := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"Name/Email/Username/Password cannot be empty."}`)
}

func (s *S) TestDeleteUser(c *C) {
	alice.Save()
	defer alice.Delete()

	s.router.Delete("/api/users", s.Api.Route(usersHandler, "DeleteUser"))
	req, _ := http.NewRequest("DELETE", "/api/users", nil)
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 200)
	c.Assert(s.recorder.Body.String(), Equals, `{"name":"Alice","email":"alice@example.org","username":"alice"}`)
}

func (s *S) TestDeleteUserWithNotSignedUser(c *C) {
	s.router.Delete("/api/users", s.Api.Route(usersHandler, "DeleteUser"))
	req, _ := http.NewRequest("DELETE", "/api/users", nil)
	s.env[CurrentUser] = "invalid-user"
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, 400)
	c.Assert(s.recorder.Body.String(), Equals, `{"status_code":400,"message":"User is not signed in."}`)
}
