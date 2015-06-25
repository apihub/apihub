package gateway

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/backstage/maestro/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestServer(c *C) {
	hostname, err := os.Hostname()
	var via string
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err == nil {
			via = fmt.Sprintf("%d.%d %s", r.ProtoMajor, r.ProtoMinor, hostname)
		}
		c.Assert(r.Header.Get("Via"), Equals, via)
		w.Write([]byte("OK"))
	}))
	defer target.Close()

	service := &ServiceHandler{service: &account.Service{Endpoint: "http://" + target.Listener.Addr().String(), Subdomain: "test", Timeout: 10, Disabled: false}}
	dispatcher := NewDispatcher(service)

	w := httptest.NewRecorder()
	w.Body = new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "test.backstage.dev", nil)
	dispatcher.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusOK)
	c.Assert(w.Body.String(), Equals, "OK")
	c.Assert(w.Header().Get("Via"), Equals, via)
}
