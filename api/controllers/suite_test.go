package controllers

import (
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	recorder *httptest.ResponseRecorder
	env      map[string]interface{}
}

var _ = Suite(&S{})
