package filter

import (
	"net/http"
)

type Filter func(*http.Request, *http.Response)
type Filters map[string]Filter

func (f Filters) Add(key string, value Filter) {
	f[key] = value
}

func (f Filters) Get(key string) Filter {
	return f[key]
}
