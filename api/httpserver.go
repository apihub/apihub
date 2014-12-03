package main

import (
	"github.com/albertoleal/backstage/api/system"
	"github.com/zenazn/goji"
)

func main() {
	var app = &system.Application{}
	app.Init()
	app.DrawRoutes()

	goji.Serve()
}
