// package storage provides in memory storage implementation, for test purposes.
package storage

import (
	"errors"
	"sync"

	"github.com/apihub/apihub"
)

type Memory struct {
	mtx      sync.RWMutex
	Services map[string]apihub.ServiceSpec
}

func New() *Memory {
	return &Memory{
		Services: make(map[string]apihub.ServiceSpec),
	}
}

func (m *Memory) UpsertService(s apihub.ServiceSpec) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.Services[s.Handle] = s
	return nil
}

func (m *Memory) FindServiceByHandle(handle string) (apihub.ServiceSpec, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if service, ok := m.Services[handle]; !ok {
		return apihub.ServiceSpec{}, errors.New("service not found")
	} else {
		return service, nil
	}
}
