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

	// Backends lists all backends.
	Backends() ([]Backend, error)

	// Addbackend adds a new backend in the list of available be's.
	AddBackend(be BackendInfo) error

	// RemoveBackend removes an existing backend from the list of available be's.
	RemoveBackend(be BackendInfo) error

	// Lookup returns the backend corresponding to the address specified.
	//
	// Errors:
	// * Backend not found for given address.
	Lookup(address string) (Backend, error)

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

// ServiceSpec specifies the params to add a new service.
type ServiceSpec struct {
	// Handle specifies the subdomain/host used to access the service.
	Handle        string        `json:"handle,omitempty"`
	Description   string        `json:"description,omitempty"`
	Disabled      bool          `json:"disabled,omitempty"`
	Documentation string        `json:"documentation,omitempty"`
	Timeout       int           `json:"timeout,omitempty"`
	Backends      []BackendInfo `json:"backend,omitempty"`
}
