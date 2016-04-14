package apihub

import "time"

//go:generate counterfeiter . Backend
type Backend interface {
	// AddService adds a new service in the pool.
	//
	// Errors:
	// * When the service handle is already taken.
	AddService(Service) error

	// RemoveService removes an existing service from the pool.
	//
	// Errors:
	// * When the service handle is not found.
	RemoveService(handle string) (Service, error)

	// Services lists all services.
	Services() ([]Service, error)

	// Lookup returns the service corresponding to the handle specified.
	//
	// Errors:
	// * Service not found for given handle.
	Lookup(handle string) (Service, error)

	// Start starts the backend.
	Start() error

	// Stop stops the backend.
	Stop() error
}

// Backend holds information about a backend.
type backend struct {
	Name             string
	Address          string
	HeartBeatAddress string
	HeartBeatTimeout time.Duration
	HeartBeatRetry   int
}
