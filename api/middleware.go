package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/albertoleal/backstage/auth"
	"github.com/albertoleal/backstage/errors"
	"github.com/zenazn/goji/web"
	. "github.com/zenazn/goji/web/middleware"
)

func AuthorizationMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authorization := r.Header.Get("Authorization")
		user, err := auth.GetUserFromToken(authorization)
		if err != nil {
			AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusUnauthorized, Message: "You do not have access to this resource."})
			return
		}
		SetCurrentUser(c, user)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ErrorHandlerMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
		key, ok := GetRequestError(c)
		if ok {
			body, _ := json.Marshal(key)
			w.WriteHeader(key.StatusCode)
			io.WriteString(w, string(body))
		}
	}

	return http.HandlerFunc(fn)
}

func RequestIdMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqId := GetReqID(*c)
		w.Header().Set("Request-Id", reqId)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFound := &errors.HTTPError{StatusCode: http.StatusNotFound, Message: "The resource you are looking for was not found."}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(notFound.StatusCode)
	body, _ := json.Marshal(notFound)
	fmt.Fprint(w, string(body))

	return
}
