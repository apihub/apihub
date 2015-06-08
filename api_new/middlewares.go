package api_new

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backstage/backstage/errors"
	"github.com/gorilla/context"
)

func authorizationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("AuthorizationMiddleware")
	next(rw, r)
}

func requestIdMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("REQID")
	next(rw, r)
}

func errorMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	next(rw, r)
	fmt.Println("ERROR")

	// key, ok := GetRequestError(c)
	// if ok {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	body, _ := json.Marshal(key)
	// 	w.WriteHeader(key.StatusCode)
	// 	io.WriteString(w, string(body))
	// 	return
	// }
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	nfe := HTTPError{ErrorType: errors.E_NOT_FOUND, ErrorDescription: "The resource requested does not exist."}
	notFound := HTTPResponse{StatusCode: http.StatusNotFound, Body: nfe}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(notFound.StatusCode)
	body, _ := json.Marshal(notFound)
	fmt.Fprint(w, string(body))
}

func contextClearerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer context.Clear(r)
	next(w, r)
}
