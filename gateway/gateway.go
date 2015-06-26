// Package gateway provide a reverse proxy with middlewares and transformers.
package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/gateway/middleware"
	"github.com/backstage/maestro/gateway/transformer"
	. "github.com/backstage/maestro/log"
)

type Settings struct {
	Host string
	Port string
}

type Gateway struct {
	pubsub       account.PubSub
	Settings     *Settings
	services     map[string]ServiceHandler
	transformers transformer.Transformers
	middlewares  middleware.Middlewares
}

func New(config *Settings, pubsub account.PubSub) *Gateway {
	g := &Gateway{
		pubsub:       pubsub,
		Settings:     config,
		services:     make(map[string]ServiceHandler),
		middlewares:  map[string]func() middleware.Middleware{},
		transformers: map[string]transformer.Transformer{},
	}

	return g
}

func (g *Gateway) Run() {
	log.Print("Starting Backstage Maestro Gateway...")
	l, err := net.Listen("tcp", g.Settings.Port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Maestro is now ready to accept connections on port %s.", g.Settings.Port)
	log.Fatal(http.Serve(l, g))
}

// handler is responsible to check if the gateway has a service to respond the request.
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subdomain := extractSubdomainFromRequest(r)
	if serviceH, ok := g.services[subdomain]; ok {
		serviceH.handler.ServeHTTP(w, r)
		return
	}

	notFound(w)
}

// LoadServices wraps and loads the services provided.
func (g *Gateway) LoadServices(services []*account.Service) {
	if services != nil {
		for _, service := range services {
			g.AddService(service)
		}
		log.Print("Services loaded.")
	}
}

func (g *Gateway) RefreshServices() {
	receiverC := make(chan interface{})
	done := make(chan bool)

	g.pubsub.Subscribe("/services", receiverC, done)

	go func() {
		for msg := range receiverC {
			if msg != nil {
				m, ok := msg.(string)
				if !ok {
					Logger.Warn("Failed to convert message to string: %+v.", msg)
					continue
				}

				mf := bytes.NewBufferString(m)
				var service account.Service
				if err := json.NewDecoder(mf).Decode(&service); err != nil {
					Logger.Warn("Failed to decode service data: %+v.", msg)
					continue
				}

				if service.Disabled {
					g.RemoveService(&service)
				} else {
					g.AddService(&service)
				}
			}
		}
	}()
}

// Add a new service that will be used for proxying requests.
func (g *Gateway) AddService(service *account.Service) {
	h := ServiceHandler{service: service}
	h.handler = newProxyHandler(h)
	g.services[h.service.Subdomain] = h
	Logger.Info("Service added on Maestro: %+v.", service)
}

// Remove an existing service from the Gateway.
func (g *Gateway) RemoveService(service *account.Service) {
	delete(g.services, service.Subdomain)
	Logger.Info("Service removed on Maestro: %+v.", service)
}

// newProxyHandler returns an instance of Dispatch, which implements http.Handler.
// It is an instance of reverse proxy that will be available to be used by Backstage Gateway.
func newProxyHandler(e ServiceHandler) http.Handler {
	if h := e.service.Endpoint; h != "" {
		return NewDispatcher(e)
	}
	return nil
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, fmt.Sprintf(`{"error":"not_found","error_description":"%s"}`, ERR_NOT_FOUND))
}
