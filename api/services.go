package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type ServicesHandler struct {
	ApiHandler
}

func (handler *ServicesHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return OK("Hello World")
}

func (handler *ServicesHandler) CreateService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	service := &Service{
		Team: c.URLParams["team"],
	}
	err = handler.parseBody(r.Body, service)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	team, err := FindTeamByAlias(service.Team, currentUser)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			return NotFound(err.Error())
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(E_BAD_REQUEST, err.Error())
		}
	}

	err = service.Save(currentUser, team)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	service, err = FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	payload, _ := json.Marshal(service)
	return Created(string(payload))
}

func (handler *ServicesHandler) UpdateService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	service := &Service{
		Team: c.URLParams["team"],
	}
	team, err := FindTeamByAlias(service.Team, currentUser)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			return NotFound(err.Error())
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(E_BAD_REQUEST, err.Error())
		}
	}

	err = handler.parseBody(r.Body, service)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	err = service.Save(currentUser, team)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	payload, _ := json.Marshal(service)
	return OK(string(payload))
}

func (handler *ServicesHandler) DeleteService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil || service.Owner != currentUser.Email {
		return NotFound(ErrServiceNotFound.Error())
	}
	err = service.Delete()
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	payload, _ := json.Marshal(service)
	return OK(string(payload))
}

func (handler *ServicesHandler) GetServiceInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil {
		return NotFound(ErrServiceNotFound.Error())
	}

	_, err = FindTeamByAlias(service.Team, currentUser)
	if err != nil {
		switch err.(type) {
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(E_BAD_REQUEST, err.Error())
		}
	}

	result, _ := json.Marshal(service)
	return OK(string(result))
}

func (handler *ServicesHandler) GetUserServices(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	services, _ := currentUser.GetServices()
	s := CollectionSerializer{Items: services, Count: len(services)}
	payload := s.Serializer()
	return OK(payload)
}
