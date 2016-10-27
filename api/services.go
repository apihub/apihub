package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/apihub/apihub"
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
		log.Info("handle-alreay-in-use", lager.Data{"handle": spec.Handle})
		s.handleError(rw, errors.New("Handle already in use."))
		return
	}

	if err := s.storage.UpsertService(spec); err != nil {
		log.Error("failed-to-store-service", err)
		s.handleError(rw, err)
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
		log.Info("failed-to-list-services")
		s.handleError(rw, err)
		return
	}

	log.Debug("services-found", lager.Data{"services": services})
	s.writeResponse(rw, response{
		StatusCode: http.StatusOK,
		Body:       services,
	})
}
