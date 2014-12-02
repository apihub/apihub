package main

import (
	"github.com/zenazn/goji"
	"github.com/albertoleal/backstage/api/system"
)

func main() {
	var app = &system.Application{}
	app.Init()
	app.DrawRoutes()

	goji.Serve()
}
