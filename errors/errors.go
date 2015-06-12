// Package errors contains common business errors to be user all over the server.
package errors

import "errors"

const (
	E_BAD_REQUEST           string = "bad_request"
	E_FORBIDDEN_REQUEST     string = "access_denied"
	E_GATEWAY_TIMEOUT       string = "gateway_timeout"
	E_INTERNAL_SERVER_ERROR string = "internal_server_error"
	E_NOT_FOUND             string = "not_found"
	E_SERVICE_UNAVAILABLE   string = "service_unavailable"
	E_UNAUTHORIZED_REQUEST  string = "unauthorized_access"
)

var (
	ErrAuthenticationFailed   = errors.New("Authentication failed.")
	ErrRemoveOwnerFromTeam    = errors.New("It is not possible to remove the owner from the team.")
	ErrUnauthorizedAccess     = errors.New("Request refused or access is not allowed.")
	ErrBadRequest             = errors.New("The request was invalid or cannot be served.")
	ErrClientNotFound         = errors.New("Client not found.")
	ErrClientNotFoundOnTeam   = errors.New("Client not found on this team.")
	ErrInvalidTokenFormat     = errors.New("Invalid token format.")
	ErrLoginRequired          = errors.New("Invalid or expired token. Please log in with your Backstage credentials.")
	ErrOnlyOwnerHasPermission = errors.New("Only the owner has permission to perform this operation.")
	ErrServiceNotFound        = errors.New("Service not found.")
	ErrTeamNotFound           = errors.New("Team not found.")
	ErrTokenNotFound          = errors.New("Token not found.")
	ErrUserNotInTeam          = errors.New("You do not belong to this team!")
	ErrNotFound               = errors.New("The resource requested does not exist.")
	ErrConfirmationPassword   = errors.New("Your new password and confirmation password do not match or are invalid.")

	ErrUserDuplicateEntry        = errors.New("Someone already has that email. Could you try another?")
	ErrUserNotFound              = errors.New("User not found.")
	ErrUserMissingRequiredFields = errors.New("Name/Email/Password cannot be empty.")

	ErrServiceMissingRequiredFields = errors.New("Endpoint/Subdomain/Team cannot be empty.")
	ErrServiceDuplicateEntry        = errors.New("There is another service with this subdomain.")

	ErrTeamMissingRequiredFields = errors.New("Name cannot be empty.")
	ErrTeamDuplicateEntry        = errors.New("Someone already has that team alias. Could you try another?")
)

// The ValidationError type indicates that any validation has failed.
type ValidationError struct {
	Payload string
}

func (err *ValidationError) Error() string {
	return err.Payload
}

// The ForbiddenError type indicates that the user does not have
// permission to perform some operation.
type ForbiddenError struct {
	Payload string
}

func (err *ForbiddenError) Error() string {
	return err.Payload
}

type NotFoundError struct {
	Payload string
}

func (err *NotFoundError) Error() string {
	return err.Payload
}

type NotFoundErrorNEW struct {
	description error
}

func NewNotFoundErrorNEW(err error) NotFoundErrorNEW {
	return NotFoundErrorNEW{description: err}
}

func (err NotFoundErrorNEW) Error() string {
	return err.description.Error()
}

type ValidationErrorNEW struct {
	description error
}

func NewValidationErrorNEW(err error) ValidationErrorNEW {
	return ValidationErrorNEW{description: err}
}

func (err ValidationErrorNEW) Error() string {
	return err.description.Error()
}

type UnauthorizedError struct {
	description error
}

func NewUnauthorizedError(err error) UnauthorizedError {
	return UnauthorizedError{description: err}
}

func (err UnauthorizedError) Error() string {
	return err.description.Error()
}

type ForbiddenErrorNEW struct {
	description error
}

func NewForbiddenErrorNEW(err error) ForbiddenErrorNEW {
	return ForbiddenErrorNEW{description: err}
}

func (err ForbiddenErrorNEW) Error() string {
	return err.description.Error()
}
