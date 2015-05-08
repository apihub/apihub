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
	. "github.com/backstage/backstage/gateway/filter"
)

type Settings struct {
	ChannelName string
	Host        string
	Port        string
}

type ServiceHandler struct {
	handler http.Handler
	service *account.Service
	filters []Filter
}

type Gateway struct {
	Settings    *Settings
	filters     Filters
	redisClient *db.RedisClient
	services    map[string]*ServiceHandler
}

func (g *Gateway) Filter() Filters {
	if g.filters == nil {
		g.filters = map[string]Filter{}
	}
	return g.filters
}

func (g *Gateway) loadFilters() {
	//g.Filter().Add("AddSecurityHeaders", AddHeaders)
	log.Print("Default filters loaded.")
}

func NewGateway(config *Settings) *Gateway {
	g := &Gateway{
		Settings:    config,
		redisClient: db.NewRedisClient(),
	}

	g.loadFilters()
	return g
}

func (g *Gateway) LoadServices(services []*account.Service) {
	s := make(map[string]*ServiceHandler)
	if services != nil {
		s = g.wrapServices(services)
	}
	g.services = s
	for _, e := range g.services {
		e.handler = createProxy(e)
	}
	log.Print("Services loaded.")
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

func (g *Gateway) wrapServices(services []*account.Service) map[string]*ServiceHandler {
	s := make(map[string]*ServiceHandler)
	for _, serv := range services {
		h := &ServiceHandler{service: serv}
		s[serv.Subdomain] = h
		for _, f := range serv.Filters {
			if filter := g.Filter().Get(f); filter != nil {
				h.filters = append(h.filters, filter)
			}
		}
	}
	return s
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

func (g *Gateway) addService(s *account.Service) {
	h := &ServiceHandler{service: s}
	h.handler = createProxy(h)
	g.services[s.Subdomain] = h
	log.Printf("New service has been added: %s -> %s.", s.Subdomain, s.Endpoint)
}

func (g *Gateway) removeService(s *account.Service) {
	delete(g.services, s.Subdomain)
	log.Printf("Service has been removed: %s -> %s.", s.Subdomain, s.Endpoint)
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
		rp := NewDispatcher(e)
		return rp.proxy
	}
	return nil
}

func notFound(w http.ResponseWriter) {
	nf := api.NotFound(ERR_NOT_FOUND)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nf.StatusCode)
	fmt.Fprintln(w, nf.Output())
}
