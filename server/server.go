package server

import (
	"net"
	"net/http"

	"github.com/apihub/apihub"
	"github.com/pivotal-golang/lager"
)

type ApihubServer struct {
	net.Listener
	logger  lager.Logger
	backend apihub.Backend

	listenAddr string
}

func New(log lager.Logger, listenAddr string, backend apihub.Backend) *ApihubServer {
	return &ApihubServer{
		backend:    backend,
		logger:     log,
		listenAddr: listenAddr,
	}
}

func (a *ApihubServer) Start() error {
	var err error

	a.logger.Info("apihub-start-server", lager.Data{"listenAddr": a.listenAddr})
	a.Listener, err = net.Listen("tcp", a.listenAddr)
	if err != nil {
		a.logger.Error("apihub-failed-starting-server", err)
		return err
	}

	if err = a.backend.Start(); err != nil {
		a.logger.Error("apihub-failed-starting-backend", err)
		return err
	}

	go http.Serve(a.Listener, a)

	return nil
}

func (a *ApihubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (a *ApihubServer) Stop() error {
	if err := a.backend.Stop(); err != nil {
		a.logger.Error("apihub-server-stop-failed", err)
		return err
	}

	return a.Listener.Close()
}
