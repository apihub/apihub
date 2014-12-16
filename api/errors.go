package api

import "errors"

var (
	ErrBadRequest           = errors.New("The request was invalid or cannot be served.")
	ErrLoginRequired        = errors.New("Invalid or expired token. Please log in with your Backstage credentials.")
	ErrServiceNotFound      = errors.New("Service not found or you're not the owner.")
	ErrAuthenticationFailed = errors.New("Authentication failed.")
)
