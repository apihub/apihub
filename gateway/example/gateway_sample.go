package main

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/gateway"
)

func main() {
	settings := &Settings{
		Host:        "backstage.example.org",
		Port:        ":8080",
		ChannelName: "services",
	}

	services := []*account.Service{&account.Service{Endpoint: "http://gohttphelloworld.appspot.com", Subdomain: "tres", Timeout: 2, Middlewares: []string{"AuthenticationMiddleware"}, Transformers: []string{"FooTransformer"}}}

	gw := NewGateway(settings)
	gw.Transformer().Add("FooTransformer", FooTransformer)
	gw.Middleware().Add("AuthenticationMiddleware", AuthenticationMiddleware)
	gw.LoadServices(services)
	gw.RefreshServices()
	gw.Run()
}

func FooTransformer(r *http.Request, w *http.Response, body *bytes.Buffer) {
	w.Header.Set("Content-Type", "text/plain")
	body.Write([]byte("Foo"))
}

func AuthenticationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth := r.Header.Get("Authorization")
	a := strings.TrimSpace(auth)
	if a == "secret" {
		next(rw, r)
		return
	}
	rw.WriteHeader(http.StatusUnauthorized)
	rw.Write([]byte("You must be logged in."))
	return
}
