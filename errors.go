package apihub

import "fmt"

type ServiceNotFoundError struct {
	Handle string
}

func (err ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service not found with the following handle: %s.", err.Handle)
}
