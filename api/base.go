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

func Created(payload string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusCreated, Payload: payload}
}

func OK(payload string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusOK, Payload: payload}
}

func BadRequest(errorType, errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusBadRequest, ErrorType: errorType, ErrorDescription: errorDescription}
}

func Forbidden(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusForbidden, ErrorType: E_FORBIDDEN_REQUEST, ErrorDescription: errorDescription}
}

func NotFound(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusNotFound, ErrorType: E_NOT_FOUND, ErrorDescription: errorDescription}
}

func Unauthorized(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusUnauthorized, ErrorType: E_UNAUTHORIZED_REQUEST, ErrorDescription: errorDescription}
}
