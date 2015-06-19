package account

import (
	"strings"

	"github.com/backstage/maestro/errors"
)

type Service struct {
	Subdomain     string   `json:"subdomain"`
	Description   string   `json:"description,omitempty"`
	Disabled      bool     `json:"disabled,omitempty"`
	Documentation string   `json:"documentation,omitempty"`
	Endpoint      string   `json:"endpoint,omitempty"`
	Transformers  []string `json:"transformers,omitempty"`
	Owner         string   `json:"owner,omitempty"`
	Team          string   `json:"team"`
	Timeout       int      `json:"timeout,omitempty"`
}

func (service *Service) Create(owner User, team Team) error {
	service.Owner = owner.Email
	service.Subdomain = strings.ToLower(service.Subdomain)
	service.Team = team.Alias

	if !service.valid() {
		return errors.NewValidationError(errors.ErrServiceMissingRequiredFields)
	}

	if service.Exists() {
		return errors.NewValidationError(errors.ErrServiceDuplicateEntry)
	}

	sendHook(newServiceEvent("service.create", service))

	return store.UpsertService(*service)
}

func (service *Service) Update() error {
	if !service.valid() {
		return errors.NewValidationError(errors.ErrServiceMissingRequiredFields)
	}

	if !service.Exists() {
		return errors.NewNotFoundError(errors.ErrServiceNotFound)
	}

	sendHook(newServiceEvent("service.update", service))
	return store.UpsertService(*service)
}

func (service *Service) Delete(owner User) error {
	if service.Owner != owner.Email {
		return errors.NewForbiddenError(errors.ErrOnlyOwnerHasPermission)
	}

	sendHook(newServiceEvent("service.delete", service))

	return store.DeleteService(*service)
}

func DeleteServicesByTeam(team Team, owner User) error {
	services, err := store.TeamServices(team)
	if err != nil {
		return err
	}
	for _, s := range services {
		s.Delete(owner)
	}
	return nil
}

func (service Service) Exists() bool {
	_, err := store.FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		return false
	}
	return true
}

func FindServiceBySubdomain(subdomain string) (*Service, error) {
	service, err := store.FindServiceBySubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (service *Service) valid() bool {
	if service.Subdomain == "" || service.Endpoint == "" || service.Team == "" {
		return false
	}
	return true
}
