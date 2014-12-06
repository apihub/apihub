package api

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/tsuru/config"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Api struct {
}

func (api *Api) Init() {
	err := config.ReadConfigFile("config.yaml")

	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err.Error())
	}
}

func (api *Api) DrawRoutes() {
	goji.Use(RequestIdMiddleware)
	goji.NotFound(NotFoundHandler)

	// Handlers
	servicesHandler := &ServicesHandler{}
	debugHandler := &DebugHandler{}
	usersHandler := &UsersHandler{}
	teamsHandler := &TeamsHandler{}

	// Public Routes
	goji.Get("/", api.Route(servicesHandler, "Index"))
	goji.Post("/api/users", api.Route(usersHandler, "CreateUser"))
	goji.Post("/api/signin", api.Route(usersHandler, "SignIn"))
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
}

func (api *Api) Route(handler interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "apilication/json")

		methodValue := reflect.ValueOf(handler).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse)
		response := method(&c, w, r)
		// Need to check if there's any error.
		_, err := GetRequestError(&c)
		if !err {
			w.WriteHeader(response.StatusCode)
			if _, exists := c.Env["Content-Type"]; exists {
				w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
			}
			io.WriteString(w, response.Payload)
		}
	}
	return fn
}
