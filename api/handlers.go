package api

import (
	"fmt"
	"net/http"
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello ApiHub!")
}

func pingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, `{"ping":"pong"}`)
}
