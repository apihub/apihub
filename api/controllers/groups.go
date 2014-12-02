package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/api/context"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

type GroupsController struct {
	ApiController
}

func (controller *GroupsController) CreateTeam(c *web.C, w http.ResponseWriter, r *http.Request) (*HTTPResponse, error) {
	owner, err := context.GetCurrentUser(c)
	if err != nil {
		erro := &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: err.Error()}
		context.AddRequestError(c, erro)
		return nil, erro
	}

	var erro *errors.HTTPError
	group := &Group{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		erro = &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "It was not possible to handle your request. Please, try again!"}
		context.AddRequestError(c, erro)
		return nil, erro
	}
	if err = json.Unmarshal(body, group); err != nil {
		erro = &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "The request was bad-formed."}
		context.AddRequestError(c, erro)
		return nil, erro
	}

	err = group.Save(owner)
	if err != nil {
		e := err.(*errors.ValidationError)
		erro = &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: e.Message}
		context.AddRequestError(c, erro)
		return nil, erro
	}
	payload, _ := json.Marshal(group)
	response := &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
	return response, nil
}
