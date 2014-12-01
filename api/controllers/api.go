package controllers

import (
	"net/http"
)

type Controller interface {
	AddHeaders()
}

type ApiController struct{}

func (controller *ApiController) AddHeaders(w http.ResponseWriter) {}