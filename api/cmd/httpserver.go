package main

import (
	. "github.com/backstage/backstage/api"
	"github.com/zenazn/goji"
)

func main() {
	var api = &Api{}
	api.Init()
	api.DrawRoutes()

	goji.Serve()
}
