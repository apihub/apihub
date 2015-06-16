package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
	"github.com/codegangsta/negroni"
	"github.com/tylerb/graceful"
)

const (
	DEFAULT_PORT    = ":8000"
	DEFAULT_TIMEOUT = 10 * time.Second
)

type Api struct {
	auth   auth.Authenticatable
	router *router
}

func NewApi(store account.Storable) *Api {
	api := &Api{router: NewRouter(), auth: auth.NewAuth(store)}
	api.Storage(store)

	api.router.NotFoundHandler(http.HandlerFunc(api.notFoundHandler))
	api.router.AddHandler(routerArguments{Path: "/", Methods: []string{"GET"}, Handler: homeHandler})

	//  Auth (login, logout, signup)
	api.router.AddHandler(routerArguments{Path: "/auth/login", Methods: []string{"POST"}, Handler: api.userLogin})
	api.router.AddHandler(routerArguments{Path: "/auth/logout", Methods: []string{"DELETE"}, Handler: api.userLogout})
	api.router.AddHandler(routerArguments{Path: "/auth/signup", Methods: []string{"POST"}, Handler: api.userSignup})
	api.router.AddHandler(routerArguments{Path: "/auth/password", Methods: []string{"PUT"}, Handler: api.userChangePassword})

	// Middlewares
	private := api.router.AddSubrouter("/api")
	api.router.AddMiddleware("/api", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.errorMiddleware),
		negroni.HandlerFunc(api.requestIdMiddleware),
		negroni.HandlerFunc(api.authorizationMiddleware),
		negroni.HandlerFunc(api.contextClearerMiddleware),
		negroni.Wrap(private),
	))

	// Users
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/users", Methods: []string{"DELETE"}, Handler: api.userDelete})

	// Teams
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams", Methods: []string{"POST"}, Handler: teamCreate})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams", Methods: []string{"GET"}, Handler: teamList})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams/{alias}", Methods: []string{"PUT"}, Handler: teamUpdate})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams/{alias}", Methods: []string{"DELETE"}, Handler: teamDelete})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams/{alias}", Methods: []string{"GET"}, Handler: teamInfo})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams/{alias}/users", Methods: []string{"PUT"}, Handler: teamAddUsers})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/teams/{alias}/users", Methods: []string{"DELETE"}, Handler: teamRemoveUsers})

	// Services
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/services", Methods: []string{"POST"}, Handler: serviceCreate})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/services", Methods: []string{"GET"}, Handler: serviceList})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/services/{subdomain}", Methods: []string{"GET"}, Handler: serviceInfo})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/services/{subdomain}", Methods: []string{"DELETE"}, Handler: serviceDelete})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/services/{subdomain}", Methods: []string{"PUT"}, Handler: serviceUpdate})

	// Apps
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/apps", Methods: []string{"POST"}, Handler: appCreate})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/apps/{client_id}", Methods: []string{"DELETE"}, Handler: appDelete})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/apps/{client_id}", Methods: []string{"GET"}, Handler: appInfo})
	api.router.AddHandler(routerArguments{PathPrefix: "/api", Path: "/apps/{client_id}", Methods: []string{"PUT"}, Handler: appUpdate})

	return api
}

func (api *Api) Handler() http.Handler {
	return api.router.Handler()
}

// Allow to override the default authentication method.
// To be compatible, it is needed to implement the Authenticatable interface.
func (api *Api) SetAuth(auth auth.Authenticatable) {
	api.auth = auth
}

// Allow to override the default storage engine.
// To be compatible, it is needed to implement the Storable interface.
func (api *Api) Storage(store account.Storable) {
	account.Storage(store)
	api.auth = auth.NewAuth(store)
}

func (api *Api) Run() {
	graceful.Run(DEFAULT_PORT, DEFAULT_TIMEOUT, api.Handler())
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello Backstage!")
}
