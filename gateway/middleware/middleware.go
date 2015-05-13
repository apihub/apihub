package middleware

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/db"
)

// Function which modify the request.
type Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// An array of Middleware with key to be used by the gateway and service.
type Middlewares map[string]Middleware

func (f Middlewares) Add(key string, value Middleware) {
	f[key] = value
}
func (f Middlewares) Get(key string) Middleware {
	return f[key]
}

// AuthorizationMiddleware authenticates the request by checking if there is
// key in Redis following the api.AuthorizationInfo struct.
func AuthenticationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth := r.Header.Get("Authorization")
	a := strings.TrimSpace(auth)
	if len(a) > 0 {
		var tokenInfo api.AuthenticationInfo
		e := api.AuthenticationInfo{}
		get(a, &tokenInfo)
		if tokenInfo != e {
			next(rw, r)
			return
		}
	}
	err := api.Unauthorized("Request refused or access is not allowed.")
	rw.WriteHeader(err.StatusCode)
	rw.Write([]byte(err.Output()))
	return
}

func get(token string, t interface{}) error {
	conn := &db.Storage{}
	return conn.GetTokenValue(token, t)
}
