package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/pivotal-golang/lager"
)

type ApihubServer struct {
	net.Listener
	logger lager.Logger

	listenAddr    string
	listenNetwork string
	router        *Router
	server        *http.Server
}

func New(log lager.Logger, listenNetwork, listenAddr string) *ApihubServer {
	s := &ApihubServer{
		logger:        log,
		listenAddr:    listenAddr,
		listenNetwork: listenNetwork,
		router:        NewRouter(),
	}

	s.router.AddHandler(RouterArguments{Path: "/", Methods: []string{"GET"}, Handler: homeHandler})

	s.server = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.router.r.ServeHTTP(w, r)
		}),
	}

	return s
}

func (a *ApihubServer) Start() error {
	var err error

	a.logger.Info("apihub-start-server", lager.Data{"listenAddr": a.listenAddr})
	a.Listener, err = net.Listen(a.listenNetwork, a.listenAddr)
	if err != nil {
		fmt.Println(err)
		a.logger.Error("apihub-failed-starting-server", err)
		return err
	}

	go a.server.Serve(a.Listener)

	return nil
}

func (a *ApihubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.Handler()
}

func (a *ApihubServer) Stop() error {
	return a.Listener.Close()
}
