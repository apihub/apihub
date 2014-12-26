// Package errors contains common business errors to be user all over the server.
package errors

import "errors"

var (
	ErrUserNotInTeam          = errors.New("You do not belong to this team!")
	ErrOnlyOwnerHasPermission = errors.New("Only the owner has permission to perform this operation.")
	ErrInvalidTokenFormat     = errors.New("Invalid token format.")
	ErrTeamNotFound           = errors.New("Team not found.")
	ErrServiceNotFound        = errors.New("Service not found.")
	ErrTokenNotFound          = errors.New("Token not found.")

	ErrAuthenticationFailed = errors.New("Authentication failed.")
	ErrBadRequest           = errors.New("The request was invalid or cannot be served.")
	ErrLoginRequired        = errors.New("Invalid or expired token. Please log in with your Backstage credentials.")
)

// The ValidationError type indicates that any validation has failed.
type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string {
	return err.Message
}

// The ForbiddenError type indicates that the user does not have
// permission to perform some operation.
type ForbiddenError struct {
	Message string
}

func (err *ForbiddenError) Error() string {
	return err.Message
}
