package gateway

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
)

const CHANNEL_NAME = "services"

type Config struct {
	Host string
	Port string
}

type ServiceHandler struct {
	service *account.Service

	handler http.Handler
}

//TODO: need to refactor this.
type Gateway struct {
	Config      *Config
	redisClient *db.RedisClient

	services map[string]*ServiceHandler
}

func NewGateway(config *Config) (*Gateway, error) {
	g := &Gateway{
		Config:      config,
		redisClient: db.NewRedisClient(),
		services:    make(map[string]*ServiceHandler),
	}
	if err := g.loadServices(); err != nil {
		return nil, err
	}
	go g.refreshServices()
	return g, nil
}

func (g *Gateway) Run() {
	l, err := net.Listen("tcp", g.Config.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.Serve(l, g))
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := g.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Not found.", http.StatusNotFound)
}

func (g *Gateway) loadServices() error {
	services := map[string]*ServiceHandler{"tres": &ServiceHandler{service: &account.Service{Endpoint: "http://localhost:3000"}}, "cinco": &ServiceHandler{service: &account.Service{Endpoint: "http://localhost:5000"}}}

	for _, e := range services {
		e.handler = createHandler(e)
		if e.handler == nil {
			log.Printf("endpoint error: %#v", e)
		}
	}
	g.services = services

	return nil
}

func (g *Gateway) refreshServices() {
	g.redisClient.Subscribe(CHANNEL_NAME)
	for {
		fmt.Printf("cli.Receive %+v\n", g.redisClient.Receive())
	}
}

func (g *Gateway) handler(r *http.Request) http.Handler {
	h := strings.TrimSpace(r.Host)
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}

	subdomain := extractSubdomain(h)
	if _, ok := g.services[subdomain]; ok {
		return g.services[subdomain].handler
	}
	return nil
}

func extractSubdomain(host string) string {
	var subdomain string
	host_parts := strings.Split(host, ".")
	if len(host_parts) > 2 {
		subdomain = host_parts[0]
	}

	return subdomain
}

func createHandler(e *ServiceHandler) http.Handler {
	if h := e.service.Endpoint; h != "" {
		rp := NewReverseProxy(e)
		return rp.proxy
	}

	return nil
}
