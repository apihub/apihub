package controllers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type ServicesController struct {
	ApiController
}

func (controller *ServicesController) Index(c *web.C, w http.ResponseWriter, r *http.Request) (*HTTPResponse, error) {
	response := &HTTPResponse{StatusCode: http.StatusOK, Payload: "Hello World"}
	return response, nil
}
