// Package api provides interfaces to interact with account through HTTP.
package api

import (
	"flag"
	"io"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/RangelReale/osin"
	. "github.com/backstage/backstage/log"
	"github.com/tsuru/config"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const API_DEFAULT_PORT string = ":8000"

type Api struct {
	oAuthServer *osin.Server
	Config      string
	Port        string
}

func (api *Api) init() {
	Logger.Info("Show time: Starting Backstage API.")
	if api.Port == "" {
		api.Port = API_DEFAULT_PORT
	}

	err := config.ReadConfigFile(api.Config)
	if err != nil {
		Logger.Error("Error reading config file: %s", err.Error())
	}
	storage := NewOAuthMongoStorage()
	api.LoadOauthServer(storage)
}

func (api *Api) LoadOauthServer(storage osin.Storage) {
	sconfig := &osin.ServerConfig{
		AuthorizationExpiration:   300,
		AccessExpiration:          3600,
		TokenType:                 "Bearer",
		AllowedAuthorizeTypes:     osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN},
		AllowedAccessTypes:        osin.AllowedAccessType{osin.AUTHORIZATION_CODE, osin.CLIENT_CREDENTIALS, osin.REFRESH_TOKEN},
		ErrorStatusCode:           400,
		AllowClientSecretInParams: false,
		AllowGetAccessRequest:     false,
	}
	api.oAuthServer = osin.NewServer(sconfig, storage)
}

// Logger() allows to replace the log mechanism.
func (api *Api) Logger(logger Log) {
	Logger = logger
}

// Log() returns the current Log mechanism.
func (api *Api) Log() Log {
	return Logger
}

// Register all the routes to be used by the API.
// There are two kind of routes: public and private.
// "Public routes" don't need to receive a valid http authorization token.
// On the other hand, "Private routes" expects to receive a valid http authorization token.
func (api *Api) drawDefaultRoutes() {
	goji.Use(RequestIdMiddleware)
	goji.NotFound(NotFoundHandler)

	// Handlers
	servicesHandler := &ServicesHandler{}
	clientsHandler := &ClientsHandler{}
	debugHandler := &DebugHandler{}
	usersHandler := &UsersHandler{}
	teamsHandler := &TeamsHandler{}
	oauthHandler := &OAuthHandler{}

	//Assets
	staticFilesLocation, err := filepath.Abs("api/views")
	if err != nil {
		Logger.Error(err.Error())
	}
	goji.Handle("/assets/*", http.FileServer(http.Dir(staticFilesLocation)))

	// Public Routes
	goji.Get("/", api.Route(servicesHandler, "Index"))
	goji.Post("/api/users", api.Route(usersHandler, "CreateUser"))
	goji.Post("/api/login", api.Route(usersHandler, "Login"))
	Logger.Info("Public routes registered.")

	//OAuth 2.0 routes
	goji.Post("/token", api.Route(oauthHandler, "Token"))
	goji.Get("/me", api.Route(oauthHandler, "Info"))
	goji.Get("/authorize", api.Route(oauthHandler, "Authorize"))
	goji.Post("/authorize", api.Route(oauthHandler, "Authorize"))
	Logger.Info("OAuth routes registered.")
	goji.Use(ErrorMiddleware)

	// Private Routes
	privateRoutes := web.New()
	goji.Handle("/api/*", privateRoutes)
	privateRoutes.Use(middleware.SubRouter)
	privateRoutes.NotFound(NotFoundHandler)
	privateRoutes.Use(AuthorizationMiddleware)
	privateRoutes.Get("/helloworld", api.Route(debugHandler, "HelloWorld"))
	privateRoutes.Delete("/users", api.Route(usersHandler, "DeleteUser"))

	privateRoutes.Post("/teams", api.Route(teamsHandler, "CreateTeam"))
	privateRoutes.Delete("/teams/:alias", api.Route(teamsHandler, "DeleteTeam"))
	privateRoutes.Get("/teams/:alias", api.Route(teamsHandler, "GetTeamInfo"))
	privateRoutes.Get("/teams", api.Route(teamsHandler, "GetUserTeams"))
	privateRoutes.Post("/teams/:alias/users", api.Route(teamsHandler, "AddUsersToTeam"))
	privateRoutes.Delete("/teams/:alias/users", api.Route(teamsHandler, "RemoveUsersFromTeam"))

	privateRoutes.Post("/teams/:team/clients", api.Route(clientsHandler, "CreateClient"))
	privateRoutes.Get("/teams/:team/clients/:id", api.Route(clientsHandler, "GetClientInfo"))
	privateRoutes.Delete("/teams/:team/clients/:id", api.Route(clientsHandler, "DeleteClient"))

	privateRoutes.Post("/services", api.Route(servicesHandler, "CreateService"))
	privateRoutes.Delete("/services/:subdomain", api.Route(servicesHandler, "DeleteService"))
	privateRoutes.Get("/services/:subdomain", api.Route(servicesHandler, "GetServiceInfo"))

	Logger.Info("Private routes registered.")
}

// Create a router based on given handler and method.
// Use reflection to find the method and execute it.
func (api *Api) Route(handler interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		c.Env["Api"] = api

		methodValue := reflect.ValueOf(handler).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse)
		response := method(&c, w, r)
		if response != nil {
			w.WriteHeader(response.StatusCode)
			if _, exists := c.Env["Content-Type"]; exists {
				w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
			} else {
				w.Header().Set("Content-Type", "application/json")
			}
			io.WriteString(w, response.Output())
			Logger.Debug("Headers: %#v. Output: %s", w.Header(), response.Output())
		}
	}
	return fn
}

func (api *Api) Start() {
	api.init()
	api.drawDefaultRoutes()
	flag.Set("bind", api.Port)
	goji.Serve()
}
