package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type ClientsHandler struct {
	ApiHandler
}

func (handler *ClientsHandler) CreateClient(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	client := &Client{}
	err = handler.parseBody(r.Body, client)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	team, err := FindTeamByAlias(client.Team, currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	err = client.Save(currentUser, team)
	if err != nil {
		return handler.handleError(err)
	}
	payload, _ := json.Marshal(client)
	return Created(string(payload))
}

func (handler *ClientsHandler) UpdateClient(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	client, err := FindClientById(c.URLParams["id"])
	if err != nil {
		return handler.handleError(err)
	}

	err = handler.parseBody(r.Body, client)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}

	team, err := FindTeamByAlias(client.Team, currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	client.Id = c.URLParams["id"]
	err = client.Save(currentUser, team)
	if err != nil {
		return handler.handleError(err)
	}
	payload, _ := json.Marshal(client)
	return OK(string(payload))
}

func (handler *ClientsHandler) DeleteClient(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	client, err := FindClientById(c.URLParams["id"])
	if err != nil || client.Owner != currentUser.Email {
		return NotFound(ErrClientNotFoundOnTeam.Error())
	}
	err = client.Delete()
	if err != nil {
		return handler.handleError(err)
	}

	payload, _ := json.Marshal(client)
	return OK(string(payload))
}

func (handler *ClientsHandler) GetClientInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	client, err := FindClientById(c.URLParams["id"])
	if err != nil {
		return NotFound(ErrClientNotFoundOnTeam.Error())
	}

	_, err = FindTeamByAlias(client.Team, currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	result, _ := json.Marshal(client)
	return OK(string(result))
}
