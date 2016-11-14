package gateway

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/braintree/manners"

	"code.cloudfoundry.org/lager"
)

type Gateway struct {
	sync.RWMutex

	server    *manners.GracefulServer
	rpCreator ReverseProxyCreator
	Services  map[string]ReverseProxy
}

func New(port string, rpCreator ReverseProxyCreator) *Gateway {
	gw := &Gateway{
		rpCreator: rpCreator,
		Services:  make(map[string]ReverseProxy, 0),
	}

	gw.server = manners.NewWithServer(&http.Server{
		Addr:           port,
		Handler:        gw,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})

	return gw
}

func (gw *Gateway) Start(logger lager.Logger) error {
	log := logger.Session("start")
	log.Info("starting", lager.Data{"addr": gw.server.Addr})

	if err := gw.server.ListenAndServe(); err != nil {
		log.Fatal("failed-to-start", err)
	}
	return nil
}

func (gw *Gateway) Stop() bool {
	return gw.server.Close()
}

func (gw *Gateway) AddService(logger lager.Logger, spec ReverseProxySpec) error {
	log := logger.Session("add-service")
	log.Info("start", lager.Data{"spec": spec})
	defer log.Info("end")

	reverseProxy, err := gw.rpCreator.Create(log, spec)
	if err != nil {
		log.Error("failed-to-create-handler", err)
		return err
	}

	gw.Lock()
	gw.Services[spec.Handle] = reverseProxy
	gw.Unlock()

	return nil
}

func (gw *Gateway) RemoveService(handle string) error {
	panic("not implemented")
}

func (gw *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	gw.RLock()
	defer gw.RUnlock()

	handle := extractSubdomainFromRequest(req)
	if reverseProxy, ok := gw.Services[handle]; ok {
		reverseProxy.ServeHTTP(rw, req)
	}

	pageNotFound(rw)
}

func extractSubdomainFromRequest(r *http.Request) string {
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

func pageNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(rw, `{"error":"not_found","error_description":"The requested resource could not be found but may be available again in the future."}`)
}
