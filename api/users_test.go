package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	s.router.Post("/api/users", s.Api.route(usersHandler, "CreateUser"))
	req, _ := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusCreated)
	c.Assert(s.recorder.Body.String(), Equals, `{"name":"Alice","email":"alice@example.org","username":"alice"}`)
}

func (s *S) TestCreateUserWithInvalidPayloadFormat(c *C) {
	payload := `"name": "Alice"`
	b := strings.NewReader(payload)

	s.router.Post("/api/users", s.Api.route(usersHandler, "CreateUser"))
	req, _ := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestCreateUserWithMissingRequiredFields(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)

	s.router.Post("/api/users", s.Api.route(usersHandler, "CreateUser"))
	req, _ := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Name/Email/Username/Password cannot be empty."}`)
}

func (s *S) TestDeleteUser(c *C) {
	alice.Save()
	defer alice.Delete()

	s.router.Delete("/api/users", s.Api.route(usersHandler, "DeleteUser"))
	req, _ := http.NewRequest("DELETE", "/api/users", nil)
	s.env[CurrentUser] = alice
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(s.recorder.Body.String(), Equals, `{"name":"Alice","email":"alice@example.org","username":"alice"}`)
}

func (s *S) TestDeleteUserWithNotSignedUser(c *C) {
	s.router.Delete("/api/users", s.Api.route(usersHandler, "DeleteUser"))
	req, _ := http.NewRequest("DELETE", "/api/users", nil)
	s.env[CurrentUser] = "invalid-user"
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestLoginUser(c *C) {
	bob.Save()
	defer bob.Delete()
	payload := `{"email":"bob@example.org", "password":"123456"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/login", s.Api.route(usersHandler, "Login"))
	req, _ := http.NewRequest("POST", "/api/login", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
}

func (s *S) TestLoginUserWithBadCredentials(c *C) {
	bob.Save()
	defer bob.Delete()
	payload := `{"email":"bob@example.org", "password":"123"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/login", s.Api.route(usersHandler, "Login"))
	req, _ := http.NewRequest("POST", "/api/login", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Authentication failed."}`)
}

func (s *S) TestLoginUserWithMalformedRequest(c *C) {
	bob.Save()
	defer bob.Delete()
	payload := `"email":"bob@example.org", "password":"123456"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/login", s.Api.route(usersHandler, "Login"))
	req, _ := http.NewRequest("POST", "/api/login", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestLogout(c *C) {
	bob.Save()
	defer bob.Delete()
	payload := `{"email":"bob@example.org", "password":"123456"}`
	b := strings.NewReader(payload)

	s.router.Post("/api/login", s.Api.route(usersHandler, "Login"))
	req, _ := http.NewRequest("POST", "/api/login", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	dec := json.NewDecoder(strings.NewReader(s.recorder.Body.String()))
	var t account.TokenInfo
	dec.Decode(&t)

	s.router.Delete("/api/logout", s.Api.route(usersHandler, "Logout"))
	req, _ = http.NewRequest("DELETE", "/api/logout", b)
	req.Header.Set("Authorization", t.Type+"  "+t.Token)
	webC = web.C{Env: s.env}
	s.recorder = httptest.NewRecorder()
	s.router.ServeHTTPC(webC, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, http.StatusNoContent)
}

func (s *S) TestChangePassword(c *C) {
	bob.Save()
	defer bob.Delete()
	payload := `{"email":"bob@example.org", "password":"123456", "new_password": "654321", "confirmation_password": "654321"}`
	b := strings.NewReader(payload)

	s.router.Put("/api/password", s.Api.route(usersHandler, "ChangePassword"))
	req, _ := http.NewRequest("PUT", "/api/password", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, http.StatusNoContent)

	payload = `{"email":"bob@example.org", "password":"654321"}`
	b = strings.NewReader(payload)
	s.router.Post("/api/login", s.Api.route(usersHandler, "Login"))
	req, _ = http.NewRequest("POST", "/api/login", b)
	req.Header.Set("Content-Type", "application/json")
	webC = web.C{Env: s.env}
	s.recorder = httptest.NewRecorder()
	s.router.ServeHTTPC(webC, s.recorder, req)
	c.Assert(s.recorder.Code, Equals, http.StatusOK)
}

func (s *S) TestChangePasswordWithInvalidConfirmation(c *C) {
	payload := `{"email":"bob@example.org", "password":"123456", "new_password": "654321", "confirmation_password": "invalid"}`
	b := strings.NewReader(payload)

	s.router.Put("/api/password", s.Api.route(usersHandler, "ChangePassword"))
	req, _ := http.NewRequest("PUT", "/api/password", b)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(s.recorder.Body.String(), Equals, `{"error":"bad_request","error_description":"Your new password and confirmation password do not match."}`)
}
