package middleware

import (
	"net/http"
	//"strings"

	//"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/db"
)

// Middleware which modify the request.
type Middleware interface {
	Configure(cfg string)
	Serve(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// An array of Middleware with key to be used by the gateway and service.
type Middlewares map[string]func() Middleware

func (f Middlewares) Add(key string, value func() Middleware) {
	f[key] = value
}
func (f Middlewares) Get(key string) func() Middleware {
	return f[key]
}

// Get Token From Redis.
func get(token string, t interface{}) error {
	conn := &db.Storage{}
	return conn.GetTokenValue(token, t)
}
