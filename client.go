package apihub

//go:generate counterfeiter . Client
type Client interface {
	// Ping pings the ApiHub server.
	//
	// Errors:
	// * Error - indicates the ApiHub API server is in a bad state.
	Ping() error

	// AddService adds a new service.
	//
	// Errors:
	// * When the host is already taken.
	AddService(ServiceSpec) (Service, error)

	// RemoveService removes an existing service.
	//
	// Errors:
	// * When the host is not found.
	RemoveService(host string) error

	// Services lists all services.
	//
	// Errors:
	// * None.
	Services() ([]Service, error)

	// FindService returns the service with the specified host.
	//
	// Errors:
	// * Service not found.
	FindService(host string) (Service, error)

	// UpdateService updates the service with the specified host.
	//
	// Errors:
	// * Service not found.
	UpdateService(string, ServiceSpec) (Service, error)
}
