package api

import (
	"net/http"

	"github.com/albertoleal/backstage/auth"
)

func authorizationMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authorization := r.Header.Get("Authorization")
	_, _, err := auth.GetToken(authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		next(w, r)
	}
}
