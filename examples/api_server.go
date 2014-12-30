package main

import (
	"net/http"

	. "github.com/backstage/backstage/api"
	"github.com/backstage/backstage/log"
	"github.com/zenazn/goji/web"
)

type ExampleHandler struct {
	ApiHandler
}

func (handler *ExampleHandler) Hi(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return OK("Hi there!")
}

func main() {
	var api = &Api{
		Config: "config.yaml",
	}
	logger := NewCustomLogger()
	logger.SetLevel(log.DEBUG)
	api.Logger(logger)
	//api.Log().Disable()

	api.Start()
}
