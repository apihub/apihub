package api

import (
	"encoding/json"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
)

type UsersHandler struct {
	ApiHandler
}

func (handler *UsersHandler) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, &user)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, "The request was bad-formed.")
	}

	err = user.Save()
	if err != nil {
		e := err.(*errors.ValidationError)
		return ResponseError(c, http.StatusBadRequest, e.Message)
	}
	user.Password = ""
	payload, _ := json.Marshal(user)
	return &HTTPResponse{StatusCode: http.StatusCreated, Payload: string(payload)}
}

func (handler *UsersHandler) DeleteUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user, err := GetCurrentUser(c)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, err.Error())
	}

	user.Delete()
	user.Password = ""
	payload, _ := json.Marshal(user)
	return &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
}

func (handler *UsersHandler) Login(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	email, password := r.FormValue("email"), r.FormValue("password")

	token, err := Login(email, password)
	if err != nil {
		return ResponseError(c, http.StatusBadRequest, "Invalid Email or Password.")
	}

	payload, _ := json.Marshal(token)
	return &HTTPResponse{StatusCode: http.StatusOK, Payload: string(payload)}
}
