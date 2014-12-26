package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type ApiHandler struct{}

func (api *ApiHandler) getCurrentUser(c *web.C) (*User, error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (api *ApiHandler) parseBody(body io.ReadCloser, r interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &r); err != nil {
		return ErrBadRequest
	}
	return nil
}

func BadRequest(message string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusBadRequest, Message: message}
}

func Created(message string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusCreated, Message: message}
}

func OK(message string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusOK, Message: message}
}

func Forbidden(message string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusForbidden, Message: message}
}

func NotFound(message string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusNotFound, Message: message}
}

func Unauthorized(message string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusUnauthorized, Message: message}
}
