package apihub

//go:generate counterfeiter . Storage
type Storage interface {
	UpsertService(ServiceSpec) error
	FindServiceByHandle(string) (ServiceSpec, error)
}
