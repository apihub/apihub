package api_new

import (
	"fmt"
	"net/http"
	"time"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/auth_new"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
)

const (
	DEFAULT_PORT    = ":8000"
	DEFAULT_TIMEOUT = 10 * time.Second
)

type Api struct {
	auth   auth_new.Authenticatable
	router *mux.Router
}

func NewApi(store func() (account_new.Storable, error)) *Api {
	account_new.NewStorable = store

	api := &Api{router: mux.NewRouter(), auth: auth_new.NewAuth()}
	api.router.HandleFunc("/", homeHandler)
	api.router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//  Auth (login, logout, signup)
	auth := api.router.PathPrefix("/auth").Subrouter()
	auth.Methods("POST").Path("/login").HandlerFunc(api.userLogin)
	auth.Methods("DELETE").Path("/logout").HandlerFunc(userLogout)
	auth.Methods("POST").Path("/signup").HandlerFunc(userSignup)
	auth.Methods("PUT").Path("/password").HandlerFunc(userChangePassword)

	//  Private Routes
	private := mux.NewRouter()

	api.router.PathPrefix("/api").Handler(negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(errorMiddleware),
		negroni.HandlerFunc(requestIdMiddleware),
		negroni.HandlerFunc(authorizationMiddleware),
		negroni.HandlerFunc(contextClearerMiddleware),
		negroni.Wrap(private),
	))

	// Users
	private.Methods("DELETE").Path("/api/users").HandlerFunc(userDelete)

	return api
}

func (api *Api) Login(email, password string) (*auth_new.ApiToken, error) {
	user, ok := api.auth.Authenticate(email, password)
	if ok {
		token, err := api.auth.CreateUserToken(user)
		if err != nil {
			Logger.Warn(err.Error())
			return nil, err
		}
		return token, nil
	}

	return nil, errors.ErrAuthenticationFailed
}

func (api *Api) GetHandler() http.Handler {
	return api.router
}

func (api *Api) SetAuth(auth auth_new.Authenticatable) {
	api.auth = auth
}

func (api *Api) Run() {
	graceful.Run(DEFAULT_PORT, DEFAULT_TIMEOUT, api.GetHandler())
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello Backstage!")
}
