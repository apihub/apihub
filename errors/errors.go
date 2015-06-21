// Package errors contains common business errors to be user all over the server.
package errors

import "errors"

const (
	E_BAD_REQUEST          string = "bad_request"
	E_FORBIDDEN_REQUEST    string = "access_denied"
	E_NOT_FOUND            string = "not_found"
	E_SERVICE_UNAVAILABLE  string = "service_unavailable"
	E_UNAUTHORIZED_REQUEST string = "unauthorized_access"
)

var (
	ErrAuthenticationFailed   = errors.New("Authentication failed.")
	ErrRemoveOwnerFromTeam    = errors.New("It is not possible to remove the owner from the team.")
	ErrUnauthorizedAccess     = errors.New("Request refused or access is not allowed.")
	ErrBadRequest             = errors.New("The request was invalid or cannot be served.")
	ErrBadResponse            = errors.New("The response was invalid or cannot be served.")
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

	ErrAppNotFound              = errors.New("App not found.")
	ErrAppMissingRequiredFields = errors.New("Name cannot be empty.")
	ErrAppDuplicateEntry        = errors.New("There is another app with this client id.")

	ErrPluginNotFound              = errors.New("Plugin Config not found.")
	ErrPluginMissingRequiredFields = errors.New("Name and Service cannot be empty.")

	ErrHookNotFound              = errors.New("Hook not found.")
	ErrHookMissingRequiredFields = errors.New("Name, Team and Events cannot be empty.")
)

type ErrorResponse struct {
	Type        string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
}

func (err ErrorResponse) Error() string {
	return err.Description
}

func NewErrorResponse(errType, description string) ErrorResponse {
	return ErrorResponse{Type: errType, Description: description}
}

type NotFoundError struct {
	description error
}

func NewNotFoundError(err error) NotFoundError {
	return NotFoundError{description: err}
}

func (err NotFoundError) Error() string {
	return err.description.Error()
}

type ValidationError struct {
	description error
}

func NewValidationError(err error) ValidationError {
	return ValidationError{description: err}
}

func (err ValidationError) Error() string {
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

type ForbiddenError struct {
	description error
}

func NewForbiddenError(err error) ForbiddenError {
	return ForbiddenError{description: err}
}

func (err ForbiddenError) Error() string {
	return err.description.Error()
}

type InvalidHostError struct {
	description error
}

func NewInvalidHostError(err error) InvalidHostError {
	return InvalidHostError{description: err}
}

func (err InvalidHostError) Error() string {
	return "You either have not selected any target or it is invalid: " + err.description.Error()
}

type RequestError struct {
	description error
}

func NewRequestError(err error) RequestError {
	return RequestError{description: err}
}

func (err RequestError) Error() string {
	return "Failed to connect to Backstage server: " + err.description.Error()
}

type ResponseError struct {
	description error
}

func NewResponseError(err error) ResponseError {
	return ResponseError{description: err}
}

func (err ResponseError) Error() string {
	return err.description.Error()
}
