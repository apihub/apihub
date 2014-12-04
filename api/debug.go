package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type DebugController struct {
	ApiController
}

func (controller *DebugController) HelloWorld(c *web.C, w http.ResponseWriter, r *http.Request) (*HTTPResponse, error) {
	c.Env["Content-Type"] = "text/plain"
	response := &HTTPResponse{StatusCode: http.StatusOK, Payload: "Hello World"}
	return response, nil
}
