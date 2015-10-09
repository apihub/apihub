package account

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apihub/apihub/errors"
	. "github.com/apihub/apihub/log"
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

	if err := service.valid(); err != nil {
		Logger.Info("Failed to create a service with invalid data: %+v.", service)
		return err
	}

	if service.Exists() {
		Logger.Info("Failed to create a service with duplicate data: %+v.", service)
		return errors.NewValidationError(errors.ErrServiceDuplicateEntry)
	}

	err := store.UpsertService(*service)
	if err == nil {
		go publishService(service)
	}
	Logger.Info("service.Create: %+v. Err: %s.", service, err)
	return err
}

func (service *Service) Update() error {
	if err := service.valid(); err != nil {
		Logger.Info("Failed to update a service with invalid data: %+v.", service)
		return err
	}

	if !service.Exists() {
		Logger.Info("Failed to update a not-found service: %+v.", service)
		return errors.NewNotFoundError(errors.ErrServiceNotFound)
	}

	err := store.UpsertService(*service)
	if err == nil {
		go publishService(service)
	}
	Logger.Info("service.Update: %+v. Err: %s.", service, err)
	return err
}

func (service *Service) Delete(owner User) error {
	if service.Owner != owner.Email {
		Logger.Warn("Only the owner has permission to delete the following service: %s.", service.Subdomain)
		return errors.NewForbiddenError(errors.ErrOnlyOwnerHasPermission)
	}

	go store.DeletePluginsByService(*service)

	err := store.DeleteService(*service)
	if err == nil {
		service.Disabled = true
		go publishService(service)
	}
	Logger.Info("service.Delete: %+v. Err: %s.", service, err)
	return err
}

func (service Service) Exists() bool {
	_, err := FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		return false
	}
	return true
}

func (service *Service) valid() error {
	if service.Subdomain == "" || service.Endpoint == "" || service.Team == "" {
		return errors.NewValidationError(errors.ErrServiceMissingRequiredFields)
	}
	return nil
}

func DeleteServicesByTeam(team Team, owner User) error {
	services, err := store.TeamServices(team)
	if err != nil {
		return err
	}
	for _, s := range services {

		s.Delete(owner)
	}

	Logger.Info("All services were excluded from the team `%s`.", team.Alias)
	return nil
}

func FindServiceBySubdomain(subdomain string) (*Service, error) {
	service, err := store.FindServiceBySubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (service *Service) asJson() []byte {
	j, _ := json.Marshal(service)
	return j
}

func publishService(service *Service) {
	name := fmt.Sprintf("/services/%s", service.Subdomain)
	pubsub.Publish(name, service.asJson())
	Logger.Info("The following service has been published: %s (subdomain) -> %s (endpoint).", service.Subdomain, service.Endpoint)
}
