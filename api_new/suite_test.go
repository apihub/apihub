package api_new_test

import (
	"net/http/httptest"
	"testing"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/account_new/mem"
	"github.com/backstage/backstage/api_new"
	. "gopkg.in/check.v1"
)

var httpClient api_new.HTTPClient

func Test(t *testing.T) { TestingT(t) }

type S struct {
	store  func() (account_new.Storable, error)
	server *httptest.Server
}

func (s *S) SetUpTest(c *C) {
	mem := mem.New()
	s.store = func() (account_new.Storable, error) {
		return mem, nil
	}

	api := api_new.NewApi(s.store)
	s.server = httptest.NewServer(api.GetHandler())
	httpClient = api_new.NewHTTPClient(s.server.URL)
}

func (s *S) TearDownSuite(c *C) {
	s.server.Close()
}

var _ = Suite(&S{})
