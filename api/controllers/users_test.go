package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/api/context"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) SetUpTest(c *C) {
	s.controller = &UsersController{}
	s.recorder = httptest.NewRecorder()
	s.env = map[string]interface{}{}
}

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
	key, _ := context.GetRequestError(&webC)
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
	key, _ := context.GetRequestError(&webC)
	body, _ := json.Marshal(key)
	c.Assert(string(body), Equals, expected)
}
