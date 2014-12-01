package system

import (
	"io"
	"net/http"
	"reflect"

	"github.com/albertoleal/backstage/api/controllers"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Application struct {
}

func (app *Application) DrawRoutes() {
	goji.NotFound(NotFoundHandler)

	// Controllers
	servicesController := &controllers.ServicesController{}
	debugController := &controllers.DebugController{}
	usersController := &controllers.UsersController{}

	// Public Routes
	goji.Get("/", app.Route(servicesController, "Index"))
	goji.Post("/api/users", app.Route(usersController, "CreateUser"))
	goji.Post("/api/signin", app.Route(usersController, "SignIn"))
	goji.Use(ErrorHandlerMiddleware)

	// Private Routes
	api := web.New()
	goji.Handle("/api/*", api)
	api.Use(middleware.SubRouter)
	api.NotFound(NotFoundHandler)
	api.Use(AuthorizationMiddleware)
	api.Use(ErrorHandlerMiddleware)
	api.Get("/helloworld", app.Route(debugController, "HelloWorld"))
}

func (app *Application) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c *web.C, w http.ResponseWriter, r *http.Request) (*controllers.HTTPResponse, error))
		response, err := method(&c, w, r)
		if err == nil {
			w.WriteHeader(response.StatusCode)
			if _, exists := c.Env["Content-Type"]; exists {
				w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
			}
			io.WriteString(w, response.Payload)
		}
	}
	return fn
}
