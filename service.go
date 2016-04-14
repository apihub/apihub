package apihub

import "time"

//go:generate counterfeiter . Service
type Service interface {
	// Handle returns the subdomain/host used to access a service.
	Handle() string

	// Start adds a service in the service pool to handle upcoming requests.
	Start() error

	// Stop stops proxying the requests.
	//
	// If kill is false, Apihub stops proxying the requests to one of the backends
	// registered.
	//
	// If kill is true, Apihub stops proxuing the requests and remove the service
	// from the service pool.
	Stop(kill bool) error

	// Info returns information about a service.
	Info() (ServiceInfo, error)

	//Addbackend adds a new backend in the list of available be's.
	AddBackend(be Backend) error

	// RemoveBackend removes an existing backend from the list of available be's.
	RemoveBackend(be Backend) error

	// Timeout waits for the duration before returning an error to the client.
	SetTimeout(time.Duration)
	Timeout() time.Duration
}

// ServiceInfo holds information about a service.
type ServiceInfo struct {
	// Either 'active' or 'stopped'.
	State string

	// Backends available to handle upcoming requests.
	Backends []Backend
}
