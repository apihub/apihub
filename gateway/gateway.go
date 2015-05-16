// Package gateway provide a reverse proxy with middlewares and transformers.
package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/api"
	"github.com/backstage/backstage/db"
	. "github.com/backstage/backstage/gateway/middleware"
	. "github.com/backstage/backstage/gateway/transformer"
)

type Settings struct {
	ChannelName string
	Host        string
	Port        string
}

// ServiceHandler registers the handler, transformers and middlewares for the given
// service.
type ServiceHandler struct {
	handler      http.Handler
	service      *account.Service
	transformers []Transformer
	middlewares  []Middleware
}

func (s *ServiceHandler) addMiddleware(m Middleware, mc *account.MiddlewareConfig) {
	marshal, err := json.Marshal(mc.Config)
	if err != nil {
		log.Printf("Wasnt possible to register middleware `%s`. Error: %s", mc.Name, err)
		return
	}
	m.Configure(string(marshal))
	s.middlewares = append(s.middlewares, m)
	log.Printf("Middleware `%s` added successfully for service `%s`.", mc.Name, s.service.Subdomain)
}

func (s *ServiceHandler) addTransformer(name string, t Transformer) {
	s.transformers = append(s.transformers, t)
	log.Printf("Transformer `%s` added successfully for service `%s`.", name, s.service.Subdomain)
}

// Gateway is a reverse proxy.
// It is possible to add custom transformers and middlewares to be used by services.
type Gateway struct {
	Settings     *Settings
	transformers Transformers
	middlewares  Middlewares
	redisClient  *db.RedisClient
	services     map[string]*ServiceHandler
}

// Transformer() returns the transformer map that will be sent to ReverseProxy.
func (g *Gateway) Transformer() Transformers {
	return g.transformers
}

// Middleware() returns the middleware map that will be sent to ReverseProxy.
func (g *Gateway) Middleware() Middlewares {
	return g.middlewares
}

// NewGateway returns a new Gateway, serving the provided services as proxy.
// The returned Gateway simply dispatch the incoming requests to services.
func NewGateway(config *Settings) *Gateway {
	g := &Gateway{
		Settings:     config,
		middlewares:  map[string]func() Middleware{},
		redisClient:  db.NewRedisClient(),
		services:     make(map[string]*ServiceHandler),
		transformers: map[string]Transformer{},
	}

	g.loadMiddlewares()
	g.loadTransformers()
	return g
}

// LoadServices wraps and loads the services provided.
func (g *Gateway) LoadServices(services []*account.Service) {
	if services != nil {
		g.wrapServices(services)
		log.Print("Services loaded.")
	}
}

func (g *Gateway) Run() {
	log.Print("Backstage Gateway starting...")
	l, err := net.Listen("tcp", g.Settings.Port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("The proxy is now ready to accept connections on port %s.", g.Settings.Port)
	log.Fatal(http.Serve(l, g))
}

func (g *Gateway) Close() {
	g.redisClient.Close()
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := g.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}
	notFound(w)
}

func (g *Gateway) HasMiddleware(name string) bool {
	if middleware := g.Middleware().Get(name); middleware != nil {
		return true
	}
	return false
}

// Creates an instance of ServiceHandler for each service for the given array.
func (g *Gateway) wrapServices(services []*account.Service) {
	for _, serv := range services {
		h := &ServiceHandler{service: serv}
		//Extract middlewares
		middlewares, err := serv.Middlewares()
		if err == nil {
			for _, mc := range middlewares {
				// Need to confirm if the middleware is registered in the Gateway.
				if middleware := g.Middleware().Get(mc.Name); middleware != nil {
					h.addMiddleware(middleware(), mc)
				}
			}
		}
		//Extract transformers
		for _, f := range serv.Transformers {
			if transformer := g.Transformer().Get(f); transformer != nil {
				h.addTransformer(f, transformer)
			}
		}

		h.handler = createProxy(h)
		g.services[h.service.Subdomain] = h
	}
}

//FIXME: need to figure out how to test this since it's running on its own goroutine.
func (g *Gateway) RefreshServices() {
	go func() {
		channel := g.Settings.ChannelName
		if channel == "" {
			log.Fatal("Missing channel name.")
		}
		if err := g.redisClient.Subscribe(channel); err != nil {
			fmt.Println("Could not connect to Redis: " + err.Error())
			os.Exit(1)
		}

		for {
			var service account.Service
			data := g.redisClient.Receive().Data
			err := json.Unmarshal([]byte(data), &service)
			if err == nil {
				if service.Disabled {
					g.removeService(&service)
				} else {
					g.addService(&service)
				}
			}
		}
	}()
}

// Add a new service that will be used for proxying requests.
func (g *Gateway) addService(s *account.Service) {
	g.LoadServices([]*account.Service{s})
	log.Printf("New service has been added: %s -> %s.", s.Subdomain, s.Endpoint)
}

// Remove a service from the Gateway.
func (g *Gateway) removeService(s *account.Service) {
	delete(g.services, s.Subdomain)
	log.Printf("Service has been removed: %s -> %s.", s.Subdomain, s.Endpoint)
}

// handler is responsible to check if the gateway has a service to respond the request.
func (g *Gateway) handler(r *http.Request) http.Handler {
	subdomain := extractSubdomain(r)
	if _, ok := g.services[subdomain]; ok {
		return g.services[subdomain].handler
	}
	return nil
}

// Load default middlewares provided by Backstage Gateway.
func (g *Gateway) loadMiddlewares() {
	g.Middleware().Add("cors", NewCorsMiddleware)
	g.Middleware().Add("authentication", NewAuthenticationMiddleware)
	log.Print("Default middlewares loaded.")
}

// Load default transformers provided by Backstage Gateway.
func (g *Gateway) loadTransformers() {
	g.Transformer().Add("ConvertXmlToJson", ConvertXmlToJson)
	g.Transformer().Add("ConvertJsonToXml", ConvertJsonToXml)
	log.Print("Default transformers loaded.")
}

// Extract the subdomain from request.
func extractSubdomain(r *http.Request) string {
	host := strings.TrimSpace(r.Host)
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}

	var subdomain string
	host_parts := strings.Split(host, ".")
	if len(host_parts) > 2 {
		subdomain = host_parts[0]
	}
	return subdomain
}

// createProxy returns an instance of Dispatch, which implements http.Handler.
// It is an instance of reverse proxy that will be available to be used by Backstage Gateway.
func createProxy(e *ServiceHandler) http.Handler {
	if h := e.service.Endpoint; h != "" {
		return NewDispatcher(e)
	}
	return nil
}

func notFound(w http.ResponseWriter) {
	nf := api.NotFound(ERR_NOT_FOUND)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nf.StatusCode)
	fmt.Fprintln(w, nf.Output())
}
