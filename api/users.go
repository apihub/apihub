package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

type UsersHandler struct {
	ApiHandler
}

func (handler *UsersHandler) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	err = user.Save()
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}
	user.Password = ""
	payload, _ := json.Marshal(user)
	return &HTTPResponse{StatusCode: http.StatusCreated, Message: string(payload)}
}

func (handler *UsersHandler) DeleteUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user, err := GetCurrentUser(c)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	user.Delete()
	user.Password = ""
	payload, _ := json.Marshal(user)
	return &HTTPResponse{StatusCode: http.StatusOK, Message: string(payload)}
}

func (handler *UsersHandler) Login(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
	}

	token, err := Login(user)
	if err != nil {
		return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: ErrAuthenticationFailed.Error()}
	}

	payload, _ := json.Marshal(token)
	return &HTTPResponse{StatusCode: http.StatusOK, Message: string(payload)}
}
