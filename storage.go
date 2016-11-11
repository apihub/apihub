package apihub

//go:generate counterfeiter . Storage

// Storage is an interface for "storage".
type Storage interface {
	UpsertService(ServiceSpec) error
	FindServiceByHandle(string) (ServiceSpec, error)
	Services() ([]ServiceSpec, error)
	RemoveService(string) error
}
