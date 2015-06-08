package api_new

import (
	"fmt"
	"net/http"
	"time"

	// "github.com/codegangsta/negroni"
	"github.com/backstage/backstage/account_new"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
	// "log"
)

const (
	DEFAULT_PORT    = ":8000"
	DEFAULT_TIMEOUT = 10 * time.Second
)

type Api struct {
	router *mux.Router
}

func NewApi(store func() (account_new.Storable, error)) *Api {
	account_new.NewStorable = store

	api := &Api{router: mux.NewRouter()}
	api.router.HandleFunc("/", homeHandler)
	api.router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//  Auth (login, logout, signup)
	auth := api.router.PathPrefix("/auth").Subrouter()
	auth.Methods("POST").Path("/login").HandlerFunc(userLogin)
	auth.Methods("DELETE").Path("/logout").HandlerFunc(userLogout)
	auth.Methods("POST").Path("/signup").HandlerFunc(userSignup)
	auth.Methods("PUT").Path("/password").HandlerFunc(userChangePassword)

	//  Private Routes
	// private := mux.NewRouter()
	// api.router.PathPrefix("/api").Handler(negroni.New(
	// 	negroni.NewRecovery(),
	// 	negroni.HandlerFunc(requestIdMiddleware),
	// 	negroni.HandlerFunc(authorizationMiddleware),
	// 	negroni.HandlerFunc(errorMiddleware),
	// 	negroni.HandlerFunc(contextClearerMiddleware),
	// 	negroni.Wrap(private),
	// ))

	// Teams
	// teams := private.Path("/teams/{alias}").Subrouter()
	// teams.Methods("GET").HandlerFunc(TeamsInfo)
	// teams.Methods("POST").HandlerFunc(TeamsCreate)

	return api
}

func (api *Api) Run() {
	graceful.Run(DEFAULT_PORT, DEFAULT_TIMEOUT, api.GetHandler())
}

func (api *Api) GetHandler() http.Handler {
	return api.router
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello Backstage!")
}
