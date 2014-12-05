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

	// Controllers
	servicesController := &ServicesController{}
	debugController := &DebugController{}
	usersController := &UsersController{}
	groupsController := &GroupsController{}

	// Public Routes
	goji.Get("/", api.Route(servicesController, "Index"))
	goji.Post("/api/users", api.Route(usersController, "CreateUser"))
	goji.Post("/api/signin", api.Route(usersController, "SignIn"))
	goji.Use(ErrorHandlerMiddleware)

	// Private Routes
	privateRoutes := web.New()
	goji.Handle("/api/*", privateRoutes)
	privateRoutes.Use(middleware.SubRouter)
	privateRoutes.NotFound(NotFoundHandler)
	privateRoutes.Use(AuthorizationMiddleware)
	privateRoutes.Get("/helloworld", api.Route(debugController, "HelloWorld"))
	privateRoutes.Delete("/users", api.Route(usersController, "DeleteUser"))

	privateRoutes.Post("/teams", api.Route(groupsController, "CreateTeam"))
	privateRoutes.Delete("/teams/:id", api.Route(groupsController, "DeleteTeam"))
	privateRoutes.Get("/teams/:id", api.Route(groupsController, "GetTeamInfo"))
	privateRoutes.Get("/teams", api.Route(groupsController, "GetUserTeams"))
	privateRoutes.Post("/teams/:id/users", api.Route(groupsController, "AddUsersToTeam"))
	privateRoutes.Delete("/teams/:id/users", api.Route(groupsController, "RemoveUsersFromTeam"))
}

func (api *Api) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "apilication/json")

		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse)
		response := method(&c, w, r)
		w.WriteHeader(response.StatusCode)
		if _, exists := c.Env["Content-Type"]; exists {
			w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
		}
		io.WriteString(w, response.Payload)
	}
	return fn
}
