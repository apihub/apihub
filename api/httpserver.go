package api

import (
	"github.com/zenazn/goji"
)

func main() {
	var api = &Api{}
	api.Init()
	api.DrawRoutes()

	goji.Serve()
}
