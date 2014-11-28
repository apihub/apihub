package controllers

import (
	. "gopkg.in/check.v1"
)

func (s *S) TestIndex(c *C) {
	c.Assert(s.apiController.IsTrue(), Equals, true)
}
