package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	server   *ApiServer
	handler  http.HandlerFunc
	recorder *httptest.ResponseRecorder
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	var err error
	s.server, err = NewApiServer()
	if err != nil {
		panic(err)
	}
}
