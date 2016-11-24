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

	logger           lager.Logger
	listenAddr       string
	listenNetwork    string
	router           *Router
	server           *http.Server
	storage          apihub.Storage
	servicePublisher apihub.ServicePublisher
}

func New(log lager.Logger, listenNetwork, listenAddr string, storage apihub.Storage, servicePublisher apihub.ServicePublisher) *ApihubServer {
	server := &ApihubServer{
		logger:           log,
		listenAddr:       listenAddr,
		listenNetwork:    listenNetwork,
		router:           NewRouter(),
		storage:          storage,
		servicePublisher: servicePublisher,
	}

	var handlers = map[Route]http.HandlerFunc{
		Home:          http.HandlerFunc(server.homeHandler),
		Ping:          http.HandlerFunc(server.pingHandler),
		AddService:    http.HandlerFunc(server.addService),
		ListServices:  http.HandlerFunc(server.listServices),
		RemoveService: http.HandlerFunc(server.removeService),
		FindService:   http.HandlerFunc(server.findService),
		UpdateService: http.HandlerFunc(server.updateService),
	}
	for route, handler := range handlers {
		server.router.AddHandler(RouterArguments{Path: Routes[route].Path, Method: Routes[route].Method, Handler: handler})
	}
	server.router.NotFoundHandler(http.HandlerFunc(server.notFoundHandler))

	server.server = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			server.router.r.ServeHTTP(w, r)
		}),
	}

	return server
}

func (a *ApihubServer) Start(keep bool) error {
	var err error

	log := a.logger.Session("start")

	log.Info("listening", lager.Data{"listenAddr": a.listenAddr})
	a.Listener, err = net.Listen(a.listenNetwork, a.listenAddr)
	if err != nil {
		fmt.Println(err)
		log.Error("failed-to-start", err)
		return err
	}

	if keep {
		log.Info("started")
		a.server.Serve(a.Listener)
		return nil
	}

	go a.server.Serve(a.Listener)
	log.Info("started")

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
