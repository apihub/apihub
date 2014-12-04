package api

import (
	"errors"

	"github.com/albertoleal/backstage/account"
	httpErr "github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

const (
	ErrRequestKey string = "RequestError"
	CurrentUser   string = "CurrentUser"
)

var ErrUserNotSigned = errors.New("User is not signed in.")

func AddRequestError(c *web.C, error *httpErr.HTTPError) {
	c.Env[ErrRequestKey] = error
}

func GetRequestError(c *web.C) (*httpErr.HTTPError, bool) {
	val, ok := c.Env[ErrRequestKey].(*httpErr.HTTPError)
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
