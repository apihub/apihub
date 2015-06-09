package api_new

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage/errors"
)

type HTTPError struct {
	ErrorType        string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type HTTPResponse struct {
	ContentType string      `json:"-"`
	StatusCode  int         `json:"-"`
	Body        interface{} `json:"body,omitempty"`
}

func (h *HTTPResponse) ToJson() []byte {
	h.ContentType = "application/json"

	r, err := json.Marshal(h.Body)
	if err != nil {
		return []byte(err.Error())
	}
	return r
}

func handleError(rw http.ResponseWriter, err error) {
	switch err.(type) {
	case *errors.NotFoundErrorNEW:
		erro := HTTPError{ErrorType: errors.E_NOT_FOUND, ErrorDescription: err.Error()}
		NotFound(rw, erro)
	case *errors.ForbiddenError:
		erro := HTTPError{ErrorType: errors.E_FORBIDDEN_REQUEST, ErrorDescription: err.Error()}
		Forbidden(rw, erro)
	case *errors.UnauthorizedError:
		erro := HTTPError{ErrorType: errors.E_UNAUTHORIZED_REQUEST, ErrorDescription: err.Error()}
		Forbidden(rw, erro)
	default:
		erro := HTTPError{ErrorType: errors.E_BAD_REQUEST, ErrorDescription: err.Error()}
		BadRequest(rw, erro)
	}
}

func BadRequest(rw http.ResponseWriter, body interface{}) {
	resp := HTTPResponse{StatusCode: http.StatusBadRequest, Body: body}
	jsonResponse(rw, resp)
}

func Created(rw http.ResponseWriter, body interface{}) {
	resp := HTTPResponse{StatusCode: http.StatusCreated, Body: body}
	jsonResponse(rw, resp)
}

func Forbidden(rw http.ResponseWriter, body interface{}) {
	resp := HTTPResponse{StatusCode: http.StatusForbidden, Body: body}
	jsonResponse(rw, resp)
}

func NotFound(rw http.ResponseWriter, body interface{}) {
	resp := HTTPResponse{StatusCode: http.StatusNotFound, Body: body}
	jsonResponse(rw, resp)
}

func Unauthorized(rw http.ResponseWriter, body interface{}) {
	resp := HTTPResponse{StatusCode: http.StatusUnauthorized, Body: body}
	jsonResponse(rw, resp)
}

func Ok(rw http.ResponseWriter, body interface{}) {
	resp := HTTPResponse{StatusCode: http.StatusOK, Body: body}
	jsonResponse(rw, resp)
}

func jsonResponse(rw http.ResponseWriter, resp HTTPResponse) {
	body := resp.ToJson()
	rw.Header().Set("Content-Type", resp.ContentType)
	rw.WriteHeader(resp.StatusCode)
	rw.Write([]byte(body))
}
