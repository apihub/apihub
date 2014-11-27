package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/albertoleal/backstage/errors"
)

type ServiceHandler struct{}

func (s *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TODO: implement this.")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TODO: create status page.")
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFound := &errors.HTTPError{StatusCode: http.StatusNotFound, Message: "The resource you are looking for was not found."}
	w.WriteHeader(notFound.StatusCode)
	body, _ := json.Marshal(notFound)
	fmt.Fprint(w, string(body))
}
