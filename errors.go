package apihub

import "fmt"

type ServiceNotFoundError struct {
	Handle string
}

func (err ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service not found with the following handle: %s.", err.Handle)
}

type InternalError struct {
	Description string
}

func NewInternalError(description string) error {
	return &InternalError{Description: description}
}

func (ie InternalError) Error() string {
	return ie.Description

}
