package gateway

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Host string
	Port string
}

type Service struct {
	Endpoint  string
	Subdomain string
	Timeout   int

	handler http.Handler
}

type Gateway struct {
	Config *Config

	services map[string]*Service
}

func NewGateway(config *Config) (*Gateway, error) {
	g := &Gateway{
		Config:   config,
		services: make(map[string]*Service),
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
	services := map[string]*Service{"tres": &Service{Endpoint: "http://localhost:3000"}, "cinco": &Service{Endpoint: "http://localhost:5000"}}

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
	for {
		if err := g.loadServices(); err != nil {
			fmt.Printf("err %+v\n", err)
		}

		time.Sleep(10 * time.Second)
		fmt.Println("Services updated")
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

func createHandler(e *Service) http.Handler {
	if h := e.Endpoint; h != "" {
		rp := NewReverseProxy(e)
		return rp.proxy
	}

	return nil
}
