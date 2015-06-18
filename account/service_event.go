package account

import (
	"time"
)

type ServiceEvent struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Sender    *User     `json:"sender"`
	Service   *Service  `json:"service"`
}

func (e *ServiceEvent) Data() *EventData {
	return &EventData{name: e.Name, team: e.Service.Team}
}

func newServiceEvent(name string, user *User, service *Service) *ServiceEvent {
	return &ServiceEvent{
		CreatedAt: time.Now().UTC(),
		Name:      name,
		Sender:    user,
		Service:   service,
	}
}
