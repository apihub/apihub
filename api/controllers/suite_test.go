package controllers

import (
  "testing"

  . "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
  apiController *ApiController
}

var _ = Suite(&S{})