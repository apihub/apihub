package middleware

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/db"
)

type Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
type Middlewares map[string]Middleware

func (f Middlewares) Add(key string, value Middleware) {
	f[key] = value
}

func (f Middlewares) Get(key string) Middleware {
	return f[key]
}

type AuthenticationInfo struct {
	ClientId  string `json:"client_id"`
	CreatedAt string `bson:"created_at" json:"created_at"`
	Expires   int    `json:"expires"`
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	UserId    string `json:"user_id"`
}

func AuthenticationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth := r.Header.Get("Authorization")
	a := strings.TrimSpace(auth)
	if len(a) > 0 {
		var tokenInfo AuthenticationInfo
		e := AuthenticationInfo{}
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
