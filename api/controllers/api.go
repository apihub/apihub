package controllers

import (
	"net/http"
)

type Controller interface {
	AddHeaders()
}

type ApiController struct{}

func (controller *ApiController) AddHeaders(w http.ResponseWriter) {}

func (controller *ApiController) IsTrue() bool {
	return true
}
