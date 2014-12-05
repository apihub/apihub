package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

type TeamsController struct {
	ApiController
}

func (controller *TeamsController) CreateTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
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
	team := &Team{}
	if err := json.Unmarshal(body, team); err != nil {
		fmt.Print("It was not possible to create a new team.")
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}

	err = team.Save(owner)
	if err != nil {
		e := err.(*errors.ValidationError)
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: e.Message}
		AddRequestError(c, erro)
		return erro
	}
	team, err = FindTeamByName(team.Name)
	if err != nil {
		e := err.(*errors.ValidationError)
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: e.Message}
		AddRequestError(c, erro)
		return erro
	}
	payload, _ := json.Marshal(team)
	response = &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
	return response
}

func (controller *TeamsController) DeleteTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	currentUser, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}
	team, err := FindTeamById(c.URLParams["id"])
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

func (controller *TeamsController) GetUserTeams(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	currentUser, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}
	teams, _ := currentUser.GetTeams()
	payload, _ := json.Marshal(teams)
	response = &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
	return response
}

func (controller *TeamsController) GetTeamInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	currentUser, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}
	team, err := FindTeamById(c.URLParams["id"])
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "Team not found."}
		AddRequestError(c, erro)
		return erro
	}
	_, ok := team.ContainsUser(currentUser)
	if !ok {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "You do not belong to this team!"}
		AddRequestError(c, erro)
		return erro
	}
	result, _ := json.Marshal(team)
	response = &HTTPResponse{StatusCode: http.StatusOK, Payload: string(result)}
	return response
}

func (controller *TeamsController) AddUsersToTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	currentUser, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}
	team, err := FindTeamById(c.URLParams["id"])
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "Team not found."}
		AddRequestError(c, erro)
		return erro
	}

	_, ok := team.ContainsUser(currentUser)
	if !ok {
		erro := &HTTPResponse{StatusCode: http.StatusForbidden, Payload: "You do not belong to this team!"}
		AddRequestError(c, erro)
		return erro
	}

	body, err := controller.getPayload(c, r)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil || payload["users"] == nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "The request was bad-formed."}
		AddRequestError(c, erro)
		return erro
	}
	var users []string
	for _, v := range payload["users"].([]interface{}) {
		switch v.(type) {
		case string:
			user := v.(string)
			users = append(users, user)
		}
	}
	team.AddUsers(users)
	result, _ := json.Marshal(team)
	response = &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(result)}
	return response
}

func (controller *TeamsController) RemoveUsersFromTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	var response *HTTPResponse
	currentUser, err := controller.getCurrentUser(c)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}
	team, err := FindTeamById(c.URLParams["id"])
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "Team not found."}
		AddRequestError(c, erro)
		return erro
	}

	_, ok := team.ContainsUser(currentUser)
	if !ok {
		erro := &HTTPResponse{StatusCode: http.StatusForbidden, Payload: "You do not belong to this team!"}
		AddRequestError(c, erro)
		return erro
	}

	body, err := controller.getPayload(c, r)
	if err != nil {
		response = &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, response)
		return response
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil || payload["users"] == nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "The request was bad-formed."}
		AddRequestError(c, erro)
		return erro
	}
	var users []string
	for _, v := range payload["users"].([]interface{}) {
		switch v.(type) {
		case string:
			user := v.(string)
			users = append(users, user)
		}
	}
	err = team.RemoveUsers(users)
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, erro)
		return erro
	}
	result, _ := json.Marshal(team)
	response = &HTTPResponse{StatusCode: http.StatusOK, Payload: string(result)}
	return response
}
