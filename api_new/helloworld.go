package api_new

import (
	"fmt"
	"net/http"
)

func helloworld(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello World!")
}
