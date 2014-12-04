package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

type UsersController struct {
	ApiController
}

func (controller *UsersController) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) (*HTTPResponse, error) {
	body, err := controller.getPayload(c, r)
	if err != nil {
		return nil, err
	}
	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		fmt.Print("It was not possible to create a new user.")
		return nil, err
	}

	err = user.Save()
	if err != nil {
		e := err.(*errors.ValidationError)
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: e.Message}
		AddRequestError(c, erro)
		return nil, err
	}
	user.Password = ""
	payload, _ := json.Marshal(user)
	response := &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
	return response, nil
}

func (controller *UsersController) DeleteUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user, err := GetCurrentUser(c)
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, erro)
		return erro
	}

	user.Delete()
	user.Password = ""
	payload, _ := json.Marshal(user)
	response := &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
	return response
}

func (controller *UsersController) SignIn(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	username, password := r.FormValue("username"), r.FormValue("password")

	token, err := SignIn(username, password)
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: "Invalid Username or Password."}
		AddRequestError(c, erro)
		return erro
	}

	payload, _ := json.Marshal(token)
	response := &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
	return response
}
