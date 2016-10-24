package api

import "net/http"

type Route int

const (
	Home Route = iota
	Ping
	AddService
)

var Routes = map[Route]RouterArguments{
	Home:       RouterArguments{Path: "/", Method: http.MethodGet},
	Ping:       RouterArguments{Path: "/ping", Method: http.MethodGet},
	AddService: RouterArguments{Path: "/services", Method: http.MethodPost},
}
