package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type router struct {
	r        *mux.Router
	notFound http.Handler
	sub      map[string]*mux.Router
}

type routerArguments struct {
	Handler    http.HandlerFunc
	Path       string
	PathPrefix string
	Methods    []string
}

func NewRouter() *router {
	return &router{
		r:   mux.NewRouter(),
		sub: make(map[string]*mux.Router),
	}
}

func (router *router) Handler() http.Handler {
	return router.r
}

func (router *router) NotFoundHandler(h http.Handler) {
	router.notFound = h
	router.r.NotFoundHandler = h
}

func (router *router) AddSubrouter(pathPrefix string) *mux.Router {
	s := mux.NewRouter()
	s.NotFoundHandler = router.notFound
	router.sub[pathPrefix] = s
	return s
}

func (router *router) Subrouter(pathPrefix string) *mux.Router {
	return router.sub[pathPrefix]
}

func (router *router) AddMiddleware(pathPrefix string, h http.Handler) {
	router.r.PathPrefix(pathPrefix).Handler(h)
}

func (router *router) AddHandler(args routerArguments) {
	var r *mux.Router

	if sub, ok := router.sub[args.PathPrefix]; ok {
		r = sub
	} else {
		r = router.r
	}

	var prefix, path string
	if args.PathPrefix != "" {
		prefix = fmt.Sprintf("/%s", strings.Trim(args.PathPrefix, "/"))
	}
	path = fmt.Sprintf("/%s", strings.Trim(args.Path, "/"))
	r.Methods(args.Methods...).Path(fmt.Sprintf("%s%s", prefix, path)).HandlerFunc(args.Handler)
}
