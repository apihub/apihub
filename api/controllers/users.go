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

type UsersController struct {
	ApiController
}

func (controller *UsersController) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) (*HTTPResponse, bool) {
	user := &User{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "It was not possible to handle your request. Please, try again!"})
		return nil, false
	}
	if err = json.Unmarshal(body, user); err != nil {
		context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "The request was bad-formed."})
		return nil, false
	}
	err = user.Save()
	if err != nil {
		e := err.(*errors.ValidationError)
		context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: e.Message})
		return nil, false
	}
	user.Password = ""
	payload, _ := json.Marshal(user)
	response := &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
	return response, true
}
