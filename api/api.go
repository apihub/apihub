package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/auth"
	. "github.com/backstage/maestro/log"
	"github.com/codegangsta/negroni"
	"github.com/tylerb/graceful"
)

const (
	DEFAULT_EVENTS_CHANNEL_LEN = 100
	DEFAULT_PORT               = ":8000"
	DEFAULT_TIMEOUT            = 10 * time.Second
)

type Api struct {
	auth   auth.Authenticatable
	store  account.Storable
	router *Router
	Events chan Event
}

func NewApi(store account.Storable, pubsub account.PubSub) *Api {
	api := &Api{router: NewRouter(), auth: auth.NewAuth(store), Events: make(chan Event, DEFAULT_EVENTS_CHANNEL_LEN)}
	api.Storage(store)
	api.PubSub(pubsub)

	api.router.NotFoundHandler(http.HandlerFunc(api.notFoundHandler))
	api.router.AddHandler(RouterArguments{Path: "/", Methods: []string{"GET"}, Handler: homeHandler})

	//  Auth (login, logout, signup)
	api.router.AddHandler(RouterArguments{Path: "/auth/login", Methods: []string{"POST"}, Handler: api.userLogin})
	api.router.AddHandler(RouterArguments{Path: "/auth/logout", Methods: []string{"DELETE"}, Handler: api.userLogout})
	api.router.AddHandler(RouterArguments{Path: "/auth/signup", Methods: []string{"POST"}, Handler: api.userSignup})
	api.router.AddHandler(RouterArguments{Path: "/auth/password", Methods: []string{"PUT"}, Handler: api.userChangePassword})

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
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/users", Methods: []string{"DELETE"}, Handler: api.userDelete})

	// Teams
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams", Methods: []string{"POST"}, Handler: authorizationRequiredHandler(api.teamCreate)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams", Methods: []string{"GET"}, Handler: authorizationRequiredHandler(api.teamList)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams/{alias}", Methods: []string{"PUT"}, Handler: authorizationRequiredHandler(api.teamUpdate)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams/{alias}", Methods: []string{"DELETE"}, Handler: authorizationRequiredHandler(api.teamDelete)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams/{alias}", Methods: []string{"GET"}, Handler: authorizationRequiredHandler(api.teamInfo)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams/{alias}/users", Methods: []string{"PUT"}, Handler: authorizationRequiredHandler(api.teamAddUsers)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/teams/{alias}/users", Methods: []string{"DELETE"}, Handler: authorizationRequiredHandler(api.teamRemoveUsers)})

	// Services
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services", Methods: []string{"POST"}, Handler: authorizationRequiredHandler(api.serviceCreate)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services", Methods: []string{"GET"}, Handler: authorizationRequiredHandler(api.serviceList)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services/{subdomain}", Methods: []string{"GET"}, Handler: authorizationRequiredHandler(api.serviceInfo)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services/{subdomain}", Methods: []string{"DELETE"}, Handler: authorizationRequiredHandler(api.serviceDelete)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services/{subdomain}", Methods: []string{"PUT"}, Handler: authorizationRequiredHandler(api.serviceUpdate)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services/{subdomain}/plugins", Methods: []string{"PUT"}, Handler: authorizationRequiredHandler(api.pluginSubsribe)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/services/{subdomain}/plugins/{plugin_name}", Methods: []string{"DELETE"}, Handler: authorizationRequiredHandler(api.pluginUnsubsribe)})

	// Apps
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/apps", Methods: []string{"POST"}, Handler: authorizationRequiredHandler(api.appCreate)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/apps/{client_id}", Methods: []string{"DELETE"}, Handler: authorizationRequiredHandler(api.appDelete)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/apps/{client_id}", Methods: []string{"GET"}, Handler: authorizationRequiredHandler(api.appInfo)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/apps/{client_id}", Methods: []string{"PUT"}, Handler: authorizationRequiredHandler(api.appUpdate)})

	// Hooks
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/hooks", Methods: []string{"PUT"}, Handler: authorizationRequiredHandler(api.hookSave)})
	api.router.AddHandler(RouterArguments{PathPrefix: "/api", Path: "/hooks/{name}", Methods: []string{"DELETE"}, Handler: authorizationRequiredHandler(api.hookDelete)})

	return api
}

func (api *Api) Handler() http.Handler {
	return api.router.Handler()
}

// This is intend to be used when loading the api only, just to connect the maestro with maestro-gateway.
// TODO: Need to improve this.
func (api *Api) AddHook(wh account.Hook) {
	api.store.UpsertHook(wh)
}

// Allow to override the default authentication method.
// To be compatible, it is needed to implement the Authenticatable interface.
func (api *Api) SetAuth(auth auth.Authenticatable) {
	api.auth = auth
}

// Allow to override the default storage engine.
// To be compatible, it is needed to implement the Storable interface.
func (api *Api) Storage(store account.Storable) {
	api.store = store
	account.Storage(store)
	api.auth = auth.NewAuth(store)
}

// Allow to override the default pubsub engine.
// To be compatible, it is needed to implement the Subscription interface.
func (api *Api) PubSub(pubsub account.PubSub) {
	account.NewPubSub(pubsub)
}

func (api *Api) Run() {
	Logger.Info(fmt.Sprintf("Backstage Maestro is now ready to accept connections on port %s.", DEFAULT_PORT))
	graceful.Run(DEFAULT_PORT, DEFAULT_TIMEOUT, api.Handler())
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello Backstage!")
}

type AuthorizationHandler func(w http.ResponseWriter, r *http.Request, user *account.User)

func authorizationRequiredHandler(ah AuthorizationHandler) http.HandlerFunc {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUser(r)
		if err != nil {
			handleError(w, err)
			return
		}
		ah(w, r, user)
	}
	return wrapper
}
