package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

type ServicesHandler struct {
	ApiHandler
}

func (handler *ServicesHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusOK, Message: "Hello World"}
}

func (handler *ServicesHandler) CreateService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	service := &Service{}
	err = handler.parseBody(r.Body, service)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	team, err := FindTeamByAlias(service.Team)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}
	err = service.Save(currentUser, team)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}
	service, err = FindServiceBySubdomain(service.Subdomain)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}
	payload, _ := json.Marshal(service)
	return &HTTPResponse{StatusCode: http.StatusCreated, Message: string(payload)}
}

func (handler *ServicesHandler) DeleteService(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil || service.Owner != currentUser.Email {
		return &HTTPResponse{StatusCode: http.StatusForbidden, Message: ErrServiceNotFound.Error()}
	}
	err = service.Delete()
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: "It was not possible to delete your service."}
	}

	payload, _ := json.Marshal(service)
	return &HTTPResponse{StatusCode: http.StatusOK, Message: string(payload)}
}

func (handler *ServicesHandler) GetServiceInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	service, err := FindServiceBySubdomain(c.URLParams["subdomain"])
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusForbidden, Message: "Service not found or you dont belong to the team responsible for it."}
	}

	team, err := FindTeamByAlias(service.Team)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: "Team not found."}
	}
	if _, err := team.ContainsUser(currentUser); err != nil {
		return &HTTPResponse{StatusCode: http.StatusForbidden, Message: err.Error()}
	}

	result, _ := json.Marshal(service)
	return &HTTPResponse{StatusCode: http.StatusOK, Message: string(result)}
}
