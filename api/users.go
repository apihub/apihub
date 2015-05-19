package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type UsersHandler struct {
	ApiHandler
}

func (handler *UsersHandler) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return handler.handleError(err)
	}
	err = user.Save()
	if err != nil {
		return handler.handleError(err)
	}
	return Created(user.ToString())
}

func (handler *UsersHandler) DeleteUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user, err := GetCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}
	auth.RevokeTokensFor(user)
	user.Delete()
	return OK(user.ToString())
}

func (handler *UsersHandler) Login(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return handler.handleError(err)
	}
	token, err := LoginAndGetToken(user)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, ErrAuthenticationFailed.Error())
	}
	payload, _ := json.Marshal(token)
	return OK(string(payload))
}

func (handler *UsersHandler) Logout(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	authorization := r.Header.Get("Authorization")
	user, err := auth.GetUserFromToken(authorization)
	if err == nil {
		auth.RevokeTokensFor(user)
	}
	return NoContent()
}

func (handler *UsersHandler) ChangePassword(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return handler.handleError(err)
	}

	if user.NewPassword != user.ConfirmationPassword {
		return BadRequest(E_BAD_REQUEST, ErrConfirmationPassword.Error())
	}
	u, err := Login(user)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, ErrAuthenticationFailed.Error())
	}

	u.Password = user.NewPassword
	err = u.Save()
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	return NoContent()
}
