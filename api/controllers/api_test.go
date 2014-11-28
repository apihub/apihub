package controllers

import (
	. "gopkg.in/check.v1"
)

func (s *S) TestIsTrue(c *C) {
	c.Assert(s.apiController.IsTrue(), Equals, true)
}
