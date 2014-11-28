package system

import (
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/zenazn/goji/web"
)

type Application struct{}

func (app *Application) Init() {
	log.Print("Starting the Api.")
}

func (app *Application) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		c.Env["Content-Type"] = "application/json"

		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c web.C, w http.ResponseWriter, r *http.Request) (string, int))

		body, code := method(c, w, r)

		switch code {
		case http.StatusOK:
			if _, exists := c.Env["Content-Type"]; exists {
				w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
			}
			io.WriteString(w, body)
		case http.StatusSeeOther, http.StatusFound:
			http.Redirect(w, r, body, code)
		}
	}
	return fn
}
