package api

import (
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type UsersHandler struct {
	Handler
}

func (handler *UsersHandler) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	if err := handler.parseBody(r.Body, user); err != nil {
		return handler.handleError(err)
	}

	if err := user.Save(); err != nil {
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
	if err := user.Delete(); err != nil {
		return handler.handleError(err)
	}
	return OK(user.ToString())
}

func (handler *UsersHandler) Login(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	if err := handler.parseBody(r.Body, user); err != nil {
		return handler.handleError(err)
	}

	token, err := LoginAndGetToken(user)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, ErrAuthenticationFailed.Error())
	}
	return OK(token.ToString())
}

func (handler *UsersHandler) Logout(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	authorization := r.Header.Get("Authorization")

	if user, err := auth.GetUserFromToken(authorization); err == nil {
		auth.RevokeTokensFor(user)
	}
	return NoContent()
}

func (handler *UsersHandler) ChangePassword(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	if err := handler.parseBody(r.Body, user); err != nil {
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
	err = u.ChangePassword()
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	return NoContent()
}
