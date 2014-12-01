package controllers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type ServicesController struct {
	ApiController
}

func (controller *ServicesController) Index(c *web.C, w http.ResponseWriter, r *http.Request) (string, int) {
	return "{\"name\": \"simple test\"}", http.StatusOK
}
