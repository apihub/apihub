package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Router struct {
	r        *mux.Router
	notFound http.Handler
}

type RouterArguments struct {
	Handler    http.HandlerFunc
	Path       string
	PathPrefix string
	Method     string
}

func NewRouter() *Router {
	return &Router{
		r: mux.NewRouter(),
	}
}

func (router *Router) Handler() http.Handler {
	return router.r
}

func (router *Router) NotFoundHandler(h http.Handler) {
	router.notFound = h
	router.r.NotFoundHandler = h
}

func (router *Router) AddHandler(args RouterArguments) {
	r := router.r
	var prefix, path string
	if args.PathPrefix != "" {
		prefix = fmt.Sprintf("/%s", strings.Trim(args.PathPrefix, "/"))
	}
	path = fmt.Sprintf("/%s", strings.Trim(args.Path, "/"))
	r.Methods(args.Method).Path(fmt.Sprintf("%s%s", prefix, path)).HandlerFunc(args.Handler)
}
