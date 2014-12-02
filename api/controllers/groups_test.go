package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/api/context"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateTeam(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()
	defer account.DeleteGroupByName("Group")

	payload := `{"name": "Group"}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	s.env[context.CurrentUser] = user
	response, erro := groupsController.CreateTeam(&web.C{Env: s.env}, s.recorder, req)
	expected := `{"name":"Group","users":["alice"],"owner":"alice"}`
	c.Assert(erro, IsNil)
	c.Assert(response.StatusCode, Equals, 201)
	c.Assert(response.Payload, Equals, expected)
}

func (s *S) TestCreateTeamWhenUserIsNotSignedIn(c *C) {
	payload := `{"name": "Group"}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	response, erro := groupsController.CreateTeam(&web.C{Env: s.env}, s.recorder, req)
	c.Assert(response, IsNil)
	er := erro.(*errors.HTTPError)
	c.Assert(er, Not(IsNil))
	c.Assert(er.StatusCode, Equals, 400)
	c.Assert(er.Message, Equals, "User is not signed in.")
}

func (s *S) TestCreateTeamWithInvalidPayloadFormat(c *C) {
	user := &account.User{Name: "Alice", Email: "alice@example.org", Username: "alice", Password: "123456"}
	user.Save()
	defer user.Delete()

	payload := `"name": "Group"`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/teams", b)
	req.Header.Set("Content-Type", "application/json")
	s.env[context.CurrentUser] = user
	c.Assert(err, IsNil)
	webC := web.C{Env: s.env}
	_, err = groupsController.CreateTeam(&webC, s.recorder, req)
	expected := `{"status_code":400,"message":"The request was bad-formed.","url":""}`
	c.Assert(err, NotNil)
	key, _ := context.GetRequestError(&webC)
	body, _ := json.Marshal(key)
	c.Assert(string(body), Equals, expected)
}
