package api

import (
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

type ApiServer struct {
	n *negroni.Negroni
	http.Handler
	mux     *mux.Router
	muxAuth *mux.Router
}

func NewApiServer() (*ApiServer, error) {
	a := &ApiServer{}
	a.n = negroni.New(negroni.NewRecovery())
	a.drawRoutes()
	return a, nil
}

func (a *ApiServer) drawRoutes() {
	a.mux = mux.NewRouter()
	a.muxAuth = mux.NewRouter()

	a.muxAuth.Handle("/services", &ServiceHandler{}).Methods("POST")

	a.mux.HandleFunc("/debug/helloworld", HelloWorldHandler)

	a.mux.PathPrefix("/").Handler(negroni.New(negroni.HandlerFunc(authorizationMiddleware), negroni.Wrap(a.muxAuth)))
	//a.mux.PathPrefix("/services").Handler(negroni.New(negroni.HandlerFunc(authorizationMiddleware), negroni.Wrap(a.muxAuth)))

	a.n.UseHandler(a.mux)
}

func (a *ApiServer) RunServer() error {
	log.Print("Starting Backstage Api Server at :2010")
	graceful.Run(":2010", 10*time.Second, a.n)
	return nil
}
