package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type ServicesHandler struct {
	ApiHandler
}

func (handler *ServicesHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	response := &HTTPResponse{StatusCode: http.StatusOK, Payload: "Hello World"}
	return response
}

func (handler *ServicesHandler) CreateService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, err.Error())
	}

	service := &Service{}
	err = handler.parseBody(r.Body, service)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, err.Error())
	}

	team, err := FindTeamByAlias(service.Team)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, err.Error())
	}
	err = service.Save(currentUser, team)
	if err != nil {
		e := err.(*errors.ValidationError)
		return ResponseError(c, http.StatusBadRequest, e.Message)
	}
	service, err = FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		e := err.(*errors.ValidationError)
		return ResponseError(c, http.StatusBadRequest, e.Message)
	}
	payload, _ := json.Marshal(service)
	return &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
}

func (handler *ServicesHandler) DeleteService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, err.Error())
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil || service.Owner != currentUser.Email {
		return ResponseError(c, http.StatusForbidden, "Service not found or you're not the owner.")
	}
	err = service.Delete()
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, "It was not possible to delete your service.")
	}

	payload, _ := json.Marshal(service)
	return &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
}

func (handler *ServicesHandler) GetServiceInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, err.Error())
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil {
		return ResponseError(c, http.StatusForbidden, "Service not found or you dont belong to the team responsible for it.")
	}

	team, err := FindTeamByAlias(service.Team)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, "Team not found.")
	}
	if _, ok := team.ContainsUser(currentUser); !ok {
		return ResponseError(c, http.StatusForbidden, "You don not have access to this.")
	}

	result, _ := json.Marshal(service)
	return &HTTPResponse{StatusCode: http.StatusOK, Payload: string(result)}
}
