package api

import "net/http"

type Route int

const (
	Home Route = iota
	Ping
	AddService
	ListServices
	RemoveService
	FindService
	UpdateService
)

var Routes = map[Route]RouterArguments{
	Home:          RouterArguments{Path: "/", Method: http.MethodGet},
	Ping:          RouterArguments{Path: "/ping", Method: http.MethodGet},
	AddService:    RouterArguments{Path: "/services", Method: http.MethodPost},
	ListServices:  RouterArguments{Path: "/services", Method: http.MethodGet},
	RemoveService: RouterArguments{Path: "/services/{host}", Method: http.MethodDelete},
	FindService:   RouterArguments{Path: "/services/{host}", Method: http.MethodGet},
	UpdateService: RouterArguments{Path: "/services/{host}", Method: http.MethodPatch},
}
