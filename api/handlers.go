package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apihub/apihub"
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello from ApiHub!")
}

func pingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, `{"ping":"pong"}`)
}

type response struct {
	StatusCode int
	Body       interface{}
}

func (s *ApihubServer) handleError(rw http.ResponseWriter, err error) {
	switch err.(type) {
	default:
		erro := apihub.ErrorResponse{Error: apihub.E_BAD_REQUEST, Description: err.Error()}
		s.writeResponse(rw, response{
			StatusCode: http.StatusBadRequest,
			Body:       erro,
		})
	}
}

func (s *ApihubServer) writeResponse(rw http.ResponseWriter, resp response) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(resp.StatusCode)
	json.NewEncoder(rw).Encode(resp.Body)
}
