package api

import (
	"encoding/json"
	"errors"
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
		s.handleError(rw, err)
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

	if err := s.storage.UpsertService(spec); err != nil {
		log.Error("failed-to-store-service", err)
		s.handleError(rw, errors.New("Failed to add new service."))
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
