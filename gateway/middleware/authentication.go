package middleware

import (
	"net/http"
	"strings"

	"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/db"
)

// AuthenticationMiddleware authenticates the request by checking if there is
// key in Redis following the api.AuthorizationInfo struct.
type AuthenticationMiddleware struct{}

func NewAuthenticationMiddleware() Middleware {
	return &AuthenticationMiddleware{}
}
func (m *AuthenticationMiddleware) Configure(cfg string) {}

func (m *AuthenticationMiddleware) ProcessRequest(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth := r.Header.Get("Authorization")
	a := strings.TrimSpace(auth)
	if len(a) > 0 {
		var tokenInfo api.AuthenticationInfo
		e := api.AuthenticationInfo{}
		get(a, &tokenInfo)
		if tokenInfo != e {
			if tokenInfo.UserId != "" {
				r.Header.Set("Backstage-User", tokenInfo.UserId)
			}
			r.Header.Set("Backstage-ClientId", tokenInfo.ClientId)
			next(rw, r)
			return
		}
	}
	err := api.Unauthorized("Request refused or access is not allowed.")
	rw.WriteHeader(err.StatusCode)
	rw.Write([]byte(err.Output()))
	return
}

// Get Token From Redis.
func get(token string, t interface{}) error {
	conn := &db.Storage{}
	return conn.GetTokenValue(token, t)
}
