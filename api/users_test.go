package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	defer func() {
		user, err := account.FindUserByUsername("alice")
		c.Assert(err, IsNil)
		err = user.Delete()
		c.Assert(err, IsNil)
	}()
	payload := `{"name": "Alice", "email": "alice@example.org", "username": "alice", "password": "123456"}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	response, erro := s.controller.CreateUser(&web.C{Env: s.env}, s.recorder, req)
	expected := `{"name":"Alice","email":"alice@example.org","username":"alice"}`
	c.Assert(erro, IsNil)
	c.Assert(response.StatusCode, Equals, 201)
	c.Assert(response.Payload, Equals, expected)

}

func (s *S) TestCreateUserWithInvalidPayloadFormat(c *C) {
	payload := `"name": "Alice"`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	webC := web.C{Env: s.env}
	_, err = s.controller.CreateUser(&webC, s.recorder, req)
	expected := `{"status_code":400,"message":"The request was bad-formed.","url":""}`
	c.Assert(err, NotNil)
	key, _ := GetRequestError(&webC)
	body, _ := json.Marshal(key)
	c.Assert(string(body), Equals, expected)
}

func (s *S) TestCreateUserWithMissingRequiredFields(c *C) {
	payload := `{}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	webC := web.C{Env: s.env}
	_, err = s.controller.CreateUser(&webC, s.recorder, req)
	expected := `{"status_code":400,"message":"Name/Email/Username/Password cannot be empty.","url":""}`
	c.Assert(err, NotNil)
	key, _ := GetRequestError(&webC)
	body, _ := json.Marshal(key)
	c.Assert(string(body), Equals, expected)
}

func (s *S) TestDeleteUser(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	req, err := http.NewRequest("DELETE", "/api/users", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = user
	response, erro := s.controller.DeleteUser(&web.C{Env: s.env}, s.recorder, req)
	expected := `{"name":"Alice","email":"alice@example.org","username":"alice"}`
	c.Assert(erro, IsNil)
	c.Assert(response.StatusCode, Equals, 200)
	c.Assert(response.Payload, Equals, expected)
}

func (s *S) TestDeleteUserWithNotSignedUser(c *C) {
	req, err := http.NewRequest("DELETE", "/api/users", nil)
	c.Assert(err, IsNil)
	s.env[CurrentUser] = "s"
	_, erro := s.controller.DeleteUser(&web.C{Env: s.env}, s.recorder, req)
	er := erro.(*errors.HTTPError)
	c.Assert(erro, NotNil)
	c.Assert(er.StatusCode, Equals, 400)
	c.Assert(er.Message, Equals, "User is not signed in.")
}
