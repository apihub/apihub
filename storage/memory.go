// package storage provides in memory storage implementation, for test purposes.
package storage

import (
	"errors"
	"sync"

	"github.com/apihub/apihub"
)

type Memory struct {
	mtx      sync.RWMutex
	services map[string]apihub.ServiceSpec
}

func New() *Memory {
	return &Memory{
		services: make(map[string]apihub.ServiceSpec),
	}
}

func (m *Memory) AddService(s apihub.ServiceSpec) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.services[s.Host]; ok {
		return errors.New("host already in use")
	}

	m.services[s.Host] = s
	return nil
}

func (m *Memory) UpdateService(s apihub.ServiceSpec) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.services[s.Host]; !ok {
		return errors.New("service not found")
	}

	m.services[s.Host] = s
	return nil
}

func (m *Memory) FindServiceByHost(host string) (apihub.ServiceSpec, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if service, ok := m.services[host]; !ok {
		return apihub.ServiceSpec{}, errors.New("service not found")
	} else {
		return service, nil
	}
}

func (m *Memory) Services() ([]apihub.ServiceSpec, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	services := []apihub.ServiceSpec{}
	for _, service := range m.services {
		services = append(services, service)
	}

	return services, nil
}

func (m *Memory) RemoveService(host string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.services[host]; !ok {
		return errors.New("service not found")
	}

	delete(m.services, host)
	return nil
}
