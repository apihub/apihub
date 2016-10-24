// Package errors contains common business errors to be user all over the server.
package requests

import "errors"

var (
	ErrBadRequest  = errors.New("The request was invalid or cannot be served.")
	ErrBadResponse = errors.New("The response was invalid or cannot be served.")
)

type InvalidHostError struct {
	description error
}

func NewInvalidHostError(err error) InvalidHostError {
	return InvalidHostError{description: err}
}

func (err InvalidHostError) Error() string {
	return err.description.Error()
}

type RequestError struct {
	description error
}

func NewRequestError(err error) RequestError {
	return RequestError{description: err}
}

func (err RequestError) Error() string {
	return err.description.Error()
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
