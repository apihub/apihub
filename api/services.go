package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type ServicesHandler struct {
	Handler
}

func (handler *ServicesHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return OK("Hello World")
}

func (handler *ServicesHandler) CreateService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}
	service := &Service{}
	err = handler.parseBody(r.Body, service)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := FindTeamByAlias(service.Team, currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	err = service.Save(currentUser, team)
	if err != nil {
		return handler.handleError(err)
	}
	service, err = FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		return handler.handleError(err)
	}
	return Created(service.ToString())
}

func (handler *ServicesHandler) UpdateService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	service := &Service{}
	err = handler.parseBody(r.Body, service)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := FindTeamByAlias(service.Team, currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	service.Subdomain = c.URLParams["subdomain"]
	err = service.Save(currentUser, team)
	if err != nil {
		return handler.handleError(err)
	}

	service, err = FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		return handler.handleError(err)
	}
	return OK(service.ToString())
}

func (handler *ServicesHandler) DeleteService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil || service.Owner != currentUser.Email {
		return NotFound(ErrServiceNotFound.Error())
	}
	err = service.Delete()
	if err != nil {
		return handler.handleError(err)
	}

	payload, _ := json.Marshal(service)
	return OK(string(payload))
}

func (handler *ServicesHandler) GetServiceInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil {
		return NotFound(ErrServiceNotFound.Error())
	}

	_, err = FindTeamByAlias(service.Team, currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	result, _ := json.Marshal(service)
	return OK(string(result))
}

func (handler *ServicesHandler) GetUserServices(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	services, _ := currentUser.GetServices()
	s := CollectionSerializer{Items: services, Count: len(services)}
	payload := s.Serializer()
	return OK(payload)
}
