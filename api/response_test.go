package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestOutputWithError(c *C) {
	err := &HTTPResponse{
		StatusCode:       http.StatusBadRequest,
		ErrorType:        "invalid_request",
		ErrorDescription: "The request is missing a required parameter.",
	}
	c.Assert(err.Output(), Equals, `{"error":"invalid_request","error_description":"The request is missing a required parameter."}`)
}

func (s *S) TestOutputWithTeam(c *C) {
	alice := &User{Name: "Alice", Username: "alice", Email: "alice@example.org", Password: "123456"}
	team := &Team{Name: "Team", Alias: "Alias", Owner: alice.Email, Users: []string{alice.Email}}
	t, _ := json.Marshal(team)
	err := &HTTPResponse{
		StatusCode: http.StatusOK,
		Payload:    string(t),
	}
	c.Assert(err.Output(), Equals, `{"name":"Team","alias":"Alias","users":["alice@example.org"],"owner":"alice@example.org"}`)
}
