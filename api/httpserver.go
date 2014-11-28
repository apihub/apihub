package main

import (
  "github.com/gorilla/context"
  "github.com/zenazn/goji"

  "github.com/albertoleal/backstage/api/controllers"
  "github.com/albertoleal/backstage/api/system"
  "github.com/zenazn/goji/web/middleware"
  "github.com/zenazn/goji/web"
)

func main() {
  var app = &system.Application{}

  goji.NotFound(system.NotFoundHandler)
  goji.Use(context.ClearHandler)

  // Controllers
  serviceController := &controllers.ServicesController{}
  debugController := &controllers.DebugController{}

  // Public Routes
  goji.Get("/", app.Route(serviceController, "Index"))

  // Private Routes
  api := web.New()
  goji.Handle("/api/*", api)
  api.Use(middleware.SubRouter)
  api.NotFound(system.NotFoundHandler)
  api.Use(system.AuthorizationMiddleware)
  api.Use(system.ErrorHandlerMiddleware)
  api.Get("/helloworld", app.Route(debugController, "HelloWorld"))

  goji.Serve()
}