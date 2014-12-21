package api

import (
	"github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

const (
	ErrRequestKey string = "RequestError"
	CurrentUser   string = "CurrentUser"
)

// Add an error in the request context.
func AddRequestError(c *web.C, error *HTTPResponse) {
	c.Env[ErrRequestKey] = error
}

// Get an error from the request context.
// Return nil and false if nothing is found.
// Otherwise, return the error.
func GetRequestError(c *web.C) (*HTTPResponse, bool) {
	val, ok := c.Env[ErrRequestKey].(*HTTPResponse)
	if !ok {
		return nil, false
	}
	return val, true
}

// Set the current user in the request context.
func SetCurrentUser(c *web.C, user interface{}) {
	user = user.(*account.User)
	c.Env[CurrentUser] = user
}

// Get the user from the request context and check if it's still valid.
func GetCurrentUser(c *web.C) (*account.User, error) {
	user, ok := c.Env[CurrentUser].(*account.User)
	if !ok || !user.Valid() {
		return nil, ErrLoginRequired
	}
	return user, nil
}
