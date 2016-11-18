package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/apihub/apihub"
	"github.com/gorilla/mux"
)

func (s *ApihubServer) addService(rw http.ResponseWriter, r *http.Request) {
	log := s.logger.Session("add-service")
	log.Debug("start")
	defer log.Debug("end")

	var spec apihub.ServiceSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		log.Error("failed-to-parse-spec", err)
		s.handleError(rw, errors.New("Failed to parse request."))
		return
	}

	if spec.Handle == "" || len(spec.Backends) == 0 {
		s.handleError(rw, errors.New("Handle and Backend cannot be empty."))
		return
	}
	if err := s.storage.AddService(spec); err != nil {
		log.Error("failed-to-store-service", err)
		s.handleError(rw, fmt.Errorf("failed to add service: '%s'", err))
		return
	}

	if !spec.Disabled {
		if err := s.servicePublisher.Publish(log, apihub.SERVICES_PREFIX, spec); err != nil {
			log.Error("failed-to-publish-service", err)
			if cleanErr := s.storage.RemoveService(spec.Handle); cleanErr != nil {
				log.Error("failed-to-remove-service", cleanErr)
			}

			s.handleError(rw, fmt.Errorf("failed to publish service: '%s'", err))
			return
		}
	}

	log.Info("service-added", lager.Data{"service": spec})
	s.writeResponse(rw, response{
		StatusCode: http.StatusCreated,
		Body:       spec,
	})
}

func (s *ApihubServer) listServices(rw http.ResponseWriter, r *http.Request) {
	log := s.logger.Session("list-services")
	log.Debug("start")
	defer log.Debug("end")

	services, err := s.storage.Services()
	if err != nil {
		log.Error("failed-to-list-services", err)
		s.handleError(rw, errors.New("Failed to retrieve service list."))
		return
	}

	collection := Collection(services, len(services))

	log.Debug("services-found", lager.Data{"services": services})
	s.writeResponse(rw, response{
		StatusCode: http.StatusOK,
		Body:       collection,
	})
}

func (s *ApihubServer) removeService(rw http.ResponseWriter, r *http.Request) {
	log := s.logger.Session("remove-service")
	log.Debug("start")
	defer log.Debug("end")

	handle := mux.Vars(r)["handle"]

	_, err := s.storage.FindServiceByHandle(handle)
	if err != nil {
		log.Error("failed-to-find-service", err, lager.Data{"handle": handle})
		s.handleError(rw, errors.New("Handle not found."))
		return
	}

	err = s.storage.RemoveService(handle)
	if err != nil {
		log.Error("failed-to-remove-service", err)
		s.handleError(rw, errors.New("Failed to remove service."))
		return
	}

	if err := s.servicePublisher.Unpublish(log, apihub.SERVICES_PREFIX, handle); err != nil {
		log.Error("failed-to-unpublish-service", err)
	}

	log.Info("service-removed", lager.Data{"handle": handle})
	s.writeResponse(rw, response{
		StatusCode: http.StatusNoContent,
	})
}

func (s *ApihubServer) findService(rw http.ResponseWriter, r *http.Request) {
	log := s.logger.Session("find-service")
	log.Debug("start")
	defer log.Debug("end")

	handle := mux.Vars(r)["handle"]

	service, err := s.storage.FindServiceByHandle(handle)
	if err != nil {
		log.Error("failed-to-find-service", err, lager.Data{"handle": handle})
		s.handleError(rw, errors.New("Failed to find service."))
		return
	}

	log.Debug("service-found", lager.Data{"service": service})
	s.writeResponse(rw, response{
		StatusCode: http.StatusOK,
		Body:       service,
	})
}

func (s *ApihubServer) updateService(rw http.ResponseWriter, r *http.Request) {
	log := s.logger.Session("update-service")
	log.Debug("start")
	defer log.Debug("end")

	handle := mux.Vars(r)["handle"]

	service, err := s.storage.FindServiceByHandle(handle)
	if err != nil {
		log.Error("failed-to-find-service", err, lager.Data{"handle": handle})
		s.handleError(rw, errors.New("Failed to find service."))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		log.Error("failed-to-parse-spec", err)
		s.handleError(rw, errors.New("Failed to parse request."))
		return
	}

	service.Handle = handle
	if err := s.storage.UpdateService(service); err != nil {
		log.Error("failed-to-store-service", err)
		s.handleError(rw, errors.New("Failed to update service."))
		return
	}

	if service.Disabled {
		if err := s.servicePublisher.Unpublish(log, apihub.SERVICES_PREFIX, service.Handle); err != nil {
			log.Error("failed-to-unpublish-service", err)
		}
	} else {
		if err := s.servicePublisher.Publish(log, apihub.SERVICES_PREFIX, service); err != nil {
			log.Error("failed-to-publish-service", err)
		}
	}

	log.Info("service-updated", lager.Data{"service": service})
	s.writeResponse(rw, response{
		StatusCode: http.StatusOK,
		Body:       service,
	})
}
