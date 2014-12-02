package controllers

import (
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var groupsController *GroupsController

type S struct {
	recorder   *httptest.ResponseRecorder
	env        map[string]interface{}
	controller *UsersController //Should remove this.
}

func (s *S) SetUpTest(c *C) {
	s.controller = &UsersController{}
	groupsController = &GroupsController{}
	s.recorder = httptest.NewRecorder()
	s.env = map[string]interface{}{}
}

var _ = Suite(&S{})
