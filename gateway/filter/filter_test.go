package filter

import (
	"net/http"

	. "gopkg.in/check.v1"
)

func (s *S) TestGetFilter(c *C) {
	c.Check(s.filters.Get("invalid"), IsNil)
}

func (s *S) TestAddFilter(c *C) {
	c.Check(s.filters.Get("AddHeader"), IsNil)
	ah := func(r *http.Request, w *http.Response) {}
	s.filters.Add("AddHeader", ah)
	c.Check(s.filters.Get("AddHeader"), NotNil)
}
