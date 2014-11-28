package system

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	router   *web.Mux
	handler  http.HandlerFunc
	recorder *httptest.ResponseRecorder
}

var _ = Suite(&S{})
