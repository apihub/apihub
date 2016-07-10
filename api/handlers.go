package api

import (
	"fmt"
	"net/http"
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello ApiHub!")
}
