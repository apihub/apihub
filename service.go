package apihub

import (
	"time"

	"code.cloudfoundry.org/lager"
)

const SERVICES_PREFIX string = "services_"

//go:generate counterfeiter . Service
//go:generate counterfeiter . ServicePublisher
//go:generate counterfeiter . ServiceSubscriber

type Service interface {
	// Host returns the subdomain/host used to access a service.
	Host() string

	// Start adds a service in the service pool to host upcoming requests.
	Start() error

	// Stop stops proxying the requests.
	Stop() error

	// Info returns information about a service.
	Info() (ServiceSpec, error)

	// Backends returns all backends in the service.
	Backends() ([]BackendInfo, error)

	// Timeout waits for the duration before returning an error to the client.
	SetTimeout(time.Duration) error
}

type ServicePublisher interface {
	Publish(logger lager.Logger, prefix string, spec ServiceSpec) error
	Unpublish(logger lager.Logger, prefix string, host string) error
}

type ServiceSubscriber interface {
	Subscribe(logger lager.Logger, prefix string, servicesCh chan ServiceSpec, stop <-chan struct{}) error
}

// ServiceInfo holds information about a service.
type ServiceSpec struct {
	// Host specifies the subdomain/host used to access the service.
	Host   string        `json:"host"`
	Disabled bool          `json:"disabled"`
	Timeout  time.Duration `json:"timeout"` // in milliseconds
	Backends []BackendInfo `json:"backends,omitempty"`
}

// Backend holds information about a backend.
type BackendInfo struct {
	Address          string `json:"address"`
	Disabled         bool   `json:"disabled"`
	HeartBeatAddress string `json:"heart_beat_address"`
	HeartBeatTimeout int    `json:"heart_beat_timeout"`
}
