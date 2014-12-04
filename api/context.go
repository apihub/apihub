package api

import (
	"errors"

	"github.com/albertoleal/backstage/account"
	"github.com/zenazn/goji/web"
)

const (
	ErrRequestKey string = "RequestError"
	CurrentUser   string = "CurrentUser"
)

var ErrUserNotSigned = errors.New("User is not signed in.")

func AddRequestError(c *web.C, error *HTTPResponse) {
	c.Env[ErrRequestKey] = error
}

func GetRequestError(c *web.C) (*HTTPResponse, bool) {
	val, ok := c.Env[ErrRequestKey].(*HTTPResponse)
	if !ok {
		return nil, false
	}
	return val, true
}

func SetCurrentUser(c *web.C, user interface{}) {
	user = user.(*account.User)
	c.Env[CurrentUser] = user
}

func GetCurrentUser(c *web.C) (*account.User, error) {
	user, ok := c.Env[CurrentUser].(*account.User)
	if !ok || !user.Valid() {
		return nil, ErrUserNotSigned
	}
	return user, nil
}
