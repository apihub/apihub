package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/backstage/backstage/auth"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func AuthorizationMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		user, err := auth.GetUserFromToken(authorization)
		if err != nil {
			AddRequestError(c, Unauthorized("Request refused or access is not allowed."))
			return
		}
		SetCurrentUser(c, user)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ErrorMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		key, ok := GetRequestError(c)
		if ok {
			body, _ := json.Marshal(key)
			w.WriteHeader(key.StatusCode)
			io.WriteString(w, string(body))
			return
		}
	}

	return http.HandlerFunc(fn)
}

func RequestIdMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(*c)
		w.Header().Set("Request-Id", reqId)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFound := NotFound("The resource requested does not exist.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(notFound.StatusCode)
	body, _ := json.Marshal(notFound)
	fmt.Fprint(w, string(body))

	return
}
