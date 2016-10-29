package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/apihub/apihub"

	"code.cloudfoundry.org/lager"
)

type ApihubServer struct {
	net.Listener
	logger lager.Logger

	listenAddr    string
	listenNetwork string
	router        *Router
	server        *http.Server
	storage       apihub.Storage
}

func New(log lager.Logger, listenNetwork, listenAddr string, storage apihub.Storage) *ApihubServer {
	s := &ApihubServer{
		logger:        log,
		listenAddr:    listenAddr,
		listenNetwork: listenNetwork,
		router:        NewRouter(),
		storage:       storage,
	}

	var handlers = map[Route]http.HandlerFunc{
		Home:          http.HandlerFunc(homeHandler),
		Ping:          http.HandlerFunc(pingHandler),
		AddService:    http.HandlerFunc(s.addService),
		ListServices:  http.HandlerFunc(s.listServices),
		RemoveService: http.HandlerFunc(s.removeService),
	}
	for route, handler := range handlers {
		s.router.AddHandler(RouterArguments{Path: Routes[route].Path, Method: Routes[route].Method, Handler: handler})
	}

	s.server = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.router.r.ServeHTTP(w, r)
		}),
	}

	return s
}

func (a *ApihubServer) Start(keep bool) error {
	var err error

	a.logger.Info("apihub-start-server", lager.Data{"listenAddr": a.listenAddr})
	a.Listener, err = net.Listen(a.listenNetwork, a.listenAddr)
	if err != nil {
		fmt.Println(err)
		a.logger.Error("apihub-failed-starting-server", err)
		return err
	}

	if keep {
		a.logger.Info("started")
		a.server.Serve(a.Listener)
		return nil
	}

	go a.server.Serve(a.Listener)
	a.logger.Info("started")

	return nil
}

func (a *ApihubServer) Handler() http.Handler {
	return a.router.Handler()
}

func (a *ApihubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler()
}

func (a *ApihubServer) Stop() error {
	return a.Listener.Close()
}
