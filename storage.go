package apihub

//go:generate counterfeiter . Storage

// Storage is an interface for "storage".
type Storage interface {
	AddService(ServiceSpec) error
	UpdateService(ServiceSpec) error
	FindServiceByHandle(string) (ServiceSpec, error)
	Services() ([]ServiceSpec, error)
	RemoveService(string) error
}
