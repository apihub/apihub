package apihub

import "time"

//go:generate counterfeiter . Backend
type Backend interface {
	Address() string

	// Returns information about a backend.
	Info() (BackendInfo, error)

	// Start starts receiving requests.
	Start() error

	// Stop stops receiving requests.
	Stop() error
}

// Backend holds information about a backend.
type BackendInfo struct {
	Name             string
	Address          string
	HeartBeatAddress string
	HeartBeatTimeout time.Duration
	HeartBeatRetry   int
}
