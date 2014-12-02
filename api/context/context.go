package context

import (
	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

const (
	ErrRequestKey string = "RequestError"
	CurrentUser   string = "CurrentUser"
)

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

func SetCurrentUser(c *web.C, username string) {
	user := &account.User{Username: username}
	//TODO: call db under the hood to retrieve all info.
	c.Env[CurrentUser] = user
}

func GetCurrentUser(c *web.C) (*account.User, bool) {
	val, ok := c.Env[CurrentUser].(*account.User)
	if !ok {
		return nil, false
	}
	return val, true
}
