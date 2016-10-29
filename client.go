package apihub

//go:generate counterfeiter . Client
type Client interface {
	// Ping pings the ApiHub server.
	//
	// Errors:
	// * InternalError - indicates the ApiHub server is in a bad state.
	Ping() error

	// AddService adds a new service.
	//
	// Errors:
	// * When the handle is already taken.
	AddService(ServiceSpec) (Service, error)

	// RemoveService removes an existing service.
	//
	// Errors:
	// * When the handle is not found.
	RemoveService(handle string) error

	// Services lists all services.
	//
	// Errors:
	// * None.
	Services() ([]Service, error)

	// FindService returns the service with the specified handle.
	//
	// Errors:
	// * Service not found.
	FindService(handle string) (Service, error)
}
