package middleware

import (
	"net/http"

	"github.com/satori/go.uuid"
)

type RequestIdMiddleware struct{}

func NewRequestIdMiddleware() *RequestIdMiddleware {
	return &RequestIdMiddleware{}
}

func (rid *RequestIdMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.Header.Set("X-Request-Id", uuid.NewV2(uuid.DomainPerson).String())
	next(rw, r)
}
