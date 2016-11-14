package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	_, err := s.storage.FindServiceByHandle(spec.Handle)
	if err == nil {
		log.Error("failed-to-find-service", err, lager.Data{"handle": spec.Handle})
		s.handleError(rw, errors.New("Handle already in use."))
		return
	}

	//FIXME: there's a race here.

	if err := s.storage.UpsertService(spec); err != nil {
		log.Error("failed-to-store-service", err)
		s.handleError(rw, fmt.Errorf("failed to add service: '%s'", err))
		return
	}

	config := apihub.ServiceConfig{
		ServiceSpec: spec,
		Time:        time.Now(),
	}
	if err := s.servicePublisher.Publish(log, config); err != nil {
		log.Error("failed-to-publish-service", err)
		// If it fails to clean up the state it's ok to just log that
		if cleanErr := s.storage.RemoveService(spec.Handle); cleanErr != nil {
			log.Error("failed-to-remove-service", cleanErr)
		}

		s.handleError(rw, fmt.Errorf("failed to publish service: '%s'", err))
		return
	}

	log.Info("service-added", lager.Data{"serviceSpec": spec})
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
	if err := s.storage.UpsertService(service); err != nil {
		log.Error("failed-to-store-service", err)
		s.handleError(rw, errors.New("Failed to update service."))
		return
	}

	log.Info("service-updated", lager.Data{"serviceSpec": service})
	s.writeResponse(rw, response{
		StatusCode: http.StatusOK,
		Body:       service,
	})
}
