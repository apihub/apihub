package api

import (
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	"github.com/gorilla/context"
)

const (
	ErrRequestKey string = "RequestError"
	CurrentUser   string = "CurrentUser"
)

// Clear removes all values stored for a given request.
func Clear(r *http.Request) {
	context.Clear(r)
}

// Add an error in the request context.
func AddRequestError(r *http.Request, err error) {
	context.Set(r, ErrRequestKey, err)
}

// Get an error from the request context.
// Return nil and false if nothing is found.
// Otherwise, return the error.
func GetRequestError(r *http.Request) (error, bool) {
	val, ok := context.GetOk(r, ErrRequestKey)
	if !ok {
		return nil, false
	}
	return val.(error), true
}

// Set the current user in the request context.
func SetCurrentUser(r *http.Request, user *account.User) {
	context.Set(r, CurrentUser, user)
}

// Get the user from the request context and check if it's still valid.
func GetCurrentUser(r *http.Request) (*account.User, error) {
	user, ok := context.GetOk(r, CurrentUser)
	if !ok {
		return nil, errors.ErrLoginRequired
	}
	return user.(*account.User), nil
}
