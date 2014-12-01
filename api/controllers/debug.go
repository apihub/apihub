package controllers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type DebugController struct {
	ApiController
}

func (controller *DebugController) HelloWorld(c *web.C, w http.ResponseWriter, r *http.Request) (string, int) {
	c.Env["Content-Type"] = "text/plain"
	return "Hello World!", http.StatusOK
}

func (controller *DebugController) AddHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Backstage-Debug", "true")
}
