package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

type GroupsController struct {
	ApiController
}

func (controller *GroupsController) CreateTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	owner, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}

	body, err := controller.getPayload(c, r)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "The request was bad-formed."}
		AddRequestError(c, response)
		return response
	}
	group := &Group{}
	if err := json.Unmarshal(body, group); err != nil {
		fmt.Print("It was not possible to create a new team.")
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}

	err = group.Save(owner)
	if err != nil {
		e := err.(*errors.ValidationError)
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: e.Message}
		AddRequestError(c, erro)
		return erro
	}
	group, err = FindGroupByName(group.Name)
	if err != nil {
		e := err.(*errors.ValidationError)
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: e.Message}
		AddRequestError(c, erro)
		return erro
	}
	payload, _ := json.Marshal(group)
	response = &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
	return response
}

func (controller *GroupsController) DeleteTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	currentUser, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}
	team, err := FindGroupById(c.URLParams["id"])
	if err != nil || team.Owner != currentUser.Username {
		response = &HTTPResponse{StatusCode: http.StatusForbidden, Payload: "Team not found or you're not the owner."}
		AddRequestError(c, response)
		return response
	}
	err = team.Delete()
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "It was not possible to delete your team."}
		AddRequestError(c, response)
		return response
	}
	team.Id = ""
	payload, _ := json.Marshal(team)
	response = &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
	return response
}
