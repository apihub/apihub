package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type DebugHandler struct {
	Handler
}

func (handler *DebugHandler) HelloWorld(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	c.Env["Content-Type"] = "text/plain"
	return OK("Hello World")
}
