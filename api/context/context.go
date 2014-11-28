 package context

import (
  "github.com/zenazn/goji/web"
  "github.com/albertoleal/backstage/errors"
)

const ErrRequestKey string = "RequestError"

func AddRequestError(c *web.C, error *errors.HTTPError) {
  c.Env[ErrRequestKey] = error
}

func GetRequestError(c *web.C) (*errors.HTTPError, bool) {
  val, ok := c.Env[ErrRequestKey].(*errors.HTTPError)
  if !ok {
    return nil, false
  }
  return val, true
}