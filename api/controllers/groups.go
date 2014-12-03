package controllers

import (
	"encoding/json"
	"fmt"
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
	owner, err := controller.getCurrentUser(c)
	if err != nil {
		return nil, err
	}

	body, err := controller.getPayload(c, r)
	if err != nil {
		return nil, err
	}
	group := &Group{}
	if err := json.Unmarshal(body, group); err != nil {
		fmt.Print("It was not possible to create a new team.")
		return nil, err
	}

	err = group.Save(owner)
	if err != nil {
		e := err.(*errors.ValidationError)
		erro := &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: e.Message}
		context.AddRequestError(c, erro)
		return nil, erro
	}
	payload, _ := json.Marshal(group)
	response := &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
	return response, nil
}
