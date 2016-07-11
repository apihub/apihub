package api

import "net/http"

type Route int

const (
	Home Route = iota
	Ping
)

var Routes = map[Route]RouterArguments{
	Home: RouterArguments{Path: "/", Method: http.MethodGet, Handler: homeHandler},
	Ping: RouterArguments{Path: "/ping", Method: http.MethodGet, Handler: pingHandler},
}
