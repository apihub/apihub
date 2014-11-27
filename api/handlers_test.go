package api

import (
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *S) TestHelloWorldHandler(c *C) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/debug/helloworld", nil)
	if err != nil {
		c.Error(err)
	}

	HelloWorldHandler(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Assert(w.Body.String(), Equals, "Hello World!")
}
