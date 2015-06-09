package main

import (
	"github.com/backstage/backstage/account/mem"
	. "github.com/backstage/backstage/api"
	"github.com/backstage/backstage/log"
	"github.com/zenazn/goji/web"
)

type HiHandler struct {
	Handler
}

func (handler *HiHandler) Index(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	return OK("Hi from custom route!")
}

func main() {
	store := mem.New()
	var api = NewApi(store)
	logger := NewCustomLogger()
	logger.SetLevel(log.DEBUG)
	api.Logger(logger)
	//api.Log().Disable()

	api.AddPrivateRoute("GET", "/hi", &HiHandler{}, "Index")

	api.Start()
}
