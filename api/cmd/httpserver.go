package main

import (
	. "github.com/albertoleal/backstage/api"
	"github.com/zenazn/goji"
)

func main() {
	var api = &Api{}
	api.Init()
	api.DrawRoutes()

	goji.Serve()
}
