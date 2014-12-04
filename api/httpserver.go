package api

import (
	"github.com/zenazn/goji"
)

func main() {
	var app = &Application{}
	app.Init()
	app.DrawRoutes()

	goji.Serve()
}
