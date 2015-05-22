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
		Port:        ":8001",
		ChannelName: "services",
	}

	one := &account.Service{Endpoint: "http://localhost:9999", Subdomain: "one", Timeout: 2}
	two := &account.Service{Endpoint: "http://localhost:8888", Subdomain: "two", Timeout: 2}
	hw := &account.Service{Endpoint: "http://gohttphelloworld.appspot.com", Subdomain: "helloworld", Timeout: 2}
	services := []*account.Service{one, two, hw}

	//confCors := &account.PluginConfig{
	//Name:    "cors",
	//Service: hw.Subdomain,
	//Config:  map[string]interface{}{"allowed_origins": []string{"www"}, "debug": true, "allowed_methods": []string{"DELETE", "PUT"}, "allow_credentials": true, "max_age": 10},
	//}
	//confCors.Save()

	gw := NewGateway(settings)
	gw.Transformer().Add("FooTransformer", FooTransformer)
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
