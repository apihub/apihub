package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/albertoleal/backstage/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestHelloWorldHandler(c *C) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/debug/helloworld", nil)
	if err != nil {
		c.Error(err)
	}

	HelloWorldHandler(w, req)
	c.Assert(w.Body.String(), Equals, "Hello World!")
}

func (s *S) TestNotFoundHandler(c *C) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/invalid-endpoint", nil)
	if err != nil {
		c.Error(err)
	}

	NotFoundHandler(w, req)
	c.Assert(w.Code, Equals, 404)
	body := &errors.HTTPError{}
	json.Unmarshal(w.Body.Bytes(), body)
	c.Assert(body.Message, Equals, "The resource you are looking for was not found.")

}
