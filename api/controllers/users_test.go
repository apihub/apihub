package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/albertoleal/backstage/account"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) SetUpTest(c *C) {
}

func (s *S) TestCreateUser(c *C) {
	controller := &UsersController{}
	payload := `{"name": "Alice", "email": "alice@example.org", "username": "alice", "password": "123456"}`
	b := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "/api/users", b)
	req.Header.Set("Content-Type", "application/json")
	c.Assert(err, IsNil)
	recorder := httptest.NewRecorder()
	env := map[string]interface{}{}
	response, ok := controller.CreateUser(&web.C{Env: env}, recorder, req)
	expected := `{"name":"Alice","email":"alice@example.org","username":"alice"}`
	c.Assert(ok, Equals, true)
	c.Assert(response.StatusCode, Equals, 201)
	c.Assert(response.Payload, Equals, expected)

	defer func() {
		user, err := account.FindUserByUsername("alice")
		c.Assert(err, IsNil)
		err = account.DeleteUser(user)
		c.Assert(err, IsNil)
	}()
}
