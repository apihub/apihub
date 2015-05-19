package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
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
		Logger.Info("Invalid payload: %s. Original Error: '%s'.", b, err.Error())
		return ErrBadRequest
	}
	return nil
}

func (api *ApiHandler) handleError(err error) *HTTPResponse {
	switch err.(type) {
	case *NotFoundError:
		return NotFound(err.Error())
	case *ForbiddenError:
		return Forbidden(err.Error())
	default:
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
}

func BadRequest(errorType, errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusBadRequest, ErrorType: errorType, ErrorDescription: errorDescription}
}

func Created(payload string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusCreated, Payload: payload}
}

func Forbidden(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusForbidden, ErrorType: E_FORBIDDEN_REQUEST, ErrorDescription: errorDescription}
}

func GatewayTimeout(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusGatewayTimeout, ErrorType: E_GATEWAY_TIMEOUT, ErrorDescription: errorDescription}
}

func InternalServerError(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusInternalServerError, ErrorType: E_INTERNAL_SERVER_ERROR, ErrorDescription: errorDescription}
}

func NoContent() *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusNoContent}
}

func NotFound(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusNotFound, ErrorType: E_NOT_FOUND, ErrorDescription: errorDescription}
}

func OK(payload string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusOK, Payload: payload}
}

func Unauthorized(errorDescription string) *HTTPResponse {
	return &HTTPResponse{StatusCode: http.StatusUnauthorized, ErrorType: E_UNAUTHORIZED_REQUEST, ErrorDescription: errorDescription}
}
