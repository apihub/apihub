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

type Controller interface{}

type ApiController struct{}

func (api *ApiController) getCurrentUser(c *web.C) (user *User, erro error) {
	user, err := context.GetCurrentUser(c)
	if err != nil {
		erro := &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: err.Error()}
		context.AddRequestError(c, erro)
		return nil, erro
	}
	return user, nil
}

func (api *ApiController) getPayload(c *web.C, r *http.Request) ([]byte, error) {
	var erro *errors.HTTPError
	var data interface{}

	defer r.Body.Close()
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		erro := &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "It was not possible to handle your request. Please, try again!"}
		context.AddRequestError(c, erro)
		return nil, erro
	}
	if err = json.Unmarshal(payload, &data); err != nil {
		erro = &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "The request was bad-formed."}
		context.AddRequestError(c, erro)
		return nil, erro
	}
	return payload, nil
}
