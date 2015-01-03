package main

import (
	"net/http"

	. "github.com/backstage/backstage/api"
	"github.com/backstage/backstage/log"
	"github.com/zenazn/goji/web"
)

type HiHandler struct {
	ApiHandler
}

func (handler *HiHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return OK("Hi from custom route!")
}


func main() {
	var config = &Config{
		FilePath: "config.yaml",
		Port: ":8000",
	}

	var api = NewApi(config)
	logger := NewCustomLogger()
	logger.SetLevel(log.DEBUG)
	api.Logger(logger)
	// api.Log().Disable()

	api.AddPrivateRoute("GET", "/hi", &HiHandler{}, "Index")

	api.Start()
}
