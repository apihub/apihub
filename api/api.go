// Package api provides interfaces to interact with account through HTTP.
package api

import (
	"flag"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/RangelReale/osin"
	. "github.com/backstage/backstage/log"
	"github.com/tsuru/config"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const API_DEFAULT_PORT string = ":8000"

type Config struct {
	FilePath string
	Port     string
}

type Api struct {
	privateRoutes *web.Mux
	oAuthServer   *osin.Server
	Config        *Config
}

func NewApi(cfg *Config) *Api {
	var api = &Api{Config: cfg}
	api.init()
	return api
}

func (api *Api) init() {
	Logger.Info("Show time: Starting Backstage API.")
	if api.Config.Port == "" {
		api.Config.Port = API_DEFAULT_PORT
	}

	api.privateRoutes = web.New()

	err := config.ReadConfigFile(api.Config.FilePath)
	if err != nil {
		Logger.Error("Error reading config file: %s", err.Error())
	}
	storage := NewOAuthMongoStorage()
	api.LoadOauthServer(storage)
}

func (api *Api) Start() {
	api.drawDefaultRoutes()
	flag.Set("bind", api.Config.Port)
	goji.Serve()
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

// AddPrivateRoute allows to add a new route. The default middlewares will be applied.
func (api *Api) AddPrivateRoute(httpMethod, path string, handler interface{}, fn string) {
	httpMethod = strings.ToLower(httpMethod)
	r := []rune(httpMethod)
	r[0] = unicode.ToUpper(r[0])
	httpMethod = string(r)
	methodValue := reflect.ValueOf(api.privateRoutes).MethodByName(httpMethod)
	methodInterface := methodValue.Interface()
	method := methodInterface.(func(pattern interface{}, handler interface{}))
	method(path, api.route(handler, fn))
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
	assetsFilesLocation, err := filepath.Abs("../api/views")
	if err != nil {
		Logger.Error(err.Error())
	}
	goji.Handle("/assets/*", http.FileServer(http.Dir(assetsFilesLocation)))

	// Public Routes
	goji.Post("/api/users", api.route(usersHandler, "CreateUser"))
	goji.Post("/api/login", api.route(usersHandler, "Login"))
	Logger.Info("Public routes registered.")

	//OAuth 2.0 routes
	goji.Post("/login/oauth/token", api.route(oauthHandler, "Token"))
	goji.Get("/me", api.route(oauthHandler, "Info"))
	goji.Get("/login/oauth/authorize", api.route(oauthHandler, "Authorize"))
	goji.Post("/login/oauth/authorize", api.route(oauthHandler, "Authorize"))
	Logger.Info("OAuth routes registered.")
	goji.Use(ErrorMiddleware)

	// Private Routes
	goji.Handle("/api/*", api.privateRoutes)
	api.privateRoutes.Use(middleware.SubRouter)
	api.privateRoutes.NotFound(NotFoundHandler)
	api.privateRoutes.Use(AuthorizationMiddleware)
	api.privateRoutes.Get("/helloworld", api.route(debugHandler, "HelloWorld"))
	api.privateRoutes.Delete("/users", api.route(usersHandler, "DeleteUser"))

	api.privateRoutes.Post("/teams", api.route(teamsHandler, "CreateTeam"))
	api.privateRoutes.Delete("/teams/:alias", api.route(teamsHandler, "DeleteTeam"))
	api.privateRoutes.Get("/teams/:alias", api.route(teamsHandler, "GetTeamInfo"))
	api.privateRoutes.Get("/teams", api.route(teamsHandler, "GetUserTeams"))
	api.privateRoutes.Post("/teams/:alias/users", api.route(teamsHandler, "AddUsersToTeam"))
	api.privateRoutes.Delete("/teams/:alias/users", api.route(teamsHandler, "RemoveUsersFromTeam"))

	api.privateRoutes.Post("/teams/:team/clients", api.route(clientsHandler, "CreateClient"))
	api.privateRoutes.Get("/teams/:team/clients/:id", api.route(clientsHandler, "GetClientInfo"))
	api.privateRoutes.Delete("/teams/:team/clients/:id", api.route(clientsHandler, "DeleteClient"))

	api.privateRoutes.Post("/teams/:team/services", api.route(servicesHandler, "CreateService"))
	api.privateRoutes.Delete("/teams/:team/services/:subdomain", api.route(servicesHandler, "DeleteService"))
	api.privateRoutes.Get("/teams/:team/services/:subdomain", api.route(servicesHandler, "GetServiceInfo"))

	Logger.Info("Private routes registered.")
}

// Create a router based on given handler and method.
// Use reflection to find the method and execute it.
func (api *Api) route(handler interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		c.Env["Api"] = api

		Logger.Debug("[REQUEST] Headers: %#v.", r.Header)
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
			Logger.Debug("[RESPONSE] Headers: %#v. Output: %s", w.Header(), response.Output())
		}
	}
	return fn
}
