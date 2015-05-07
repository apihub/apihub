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

type Config struct {
	ChannelName string
	Host        string
	Port        string
}

type ServiceHandler struct {
	handler http.Handler
	service *account.Service
}

type Gateway struct {
	Config      *Config
	redisClient *db.RedisClient
	services    map[string]*ServiceHandler
}

func NewGateway(config *Config, services []*account.Service) *Gateway {
	s := make(map[string]*ServiceHandler)
	if services != nil {
		s = wrapService(services)
	}

	g := &Gateway{
		Config:      config,
		redisClient: db.NewRedisClient(),
		services:    s,
	}

	g.loadServices()
	return g
}

func (g *Gateway) Run() {
	l, err := net.Listen("tcp", g.Config.Port)
	if err != nil {
		log.Fatal(err)
	}

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

	http.Error(w, "Not found.", http.StatusNotFound)
}

func (g *Gateway) loadServices() {
	for _, e := range g.services {
		e.handler = createProxy(e)
	}
}

func wrapService(services []*account.Service) map[string]*ServiceHandler {
	s := make(map[string]*ServiceHandler)
	for _, serv := range services {
		s[serv.Subdomain] = &ServiceHandler{service: serv}
	}
	return s
}

func (g *Gateway) RefreshServices() {
	channel := g.Config.ChannelName
	if channel == "" {
		log.Fatal("Missing channel name.")
	}
	g.redisClient.Subscribe(channel)
	fmt.Printf("channel %+v\n", channel)
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

func createProxy(e *ServiceHandler) http.Handler {
	if h := e.service.Endpoint; h != "" {
		rp := NewReverseProxy(e)
		return rp.proxy
	}
	return nil
}
