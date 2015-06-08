package api_new

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
)

func userSignup(rw http.ResponseWriter, r *http.Request) {
	user := account_new.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	if err := user.Create(); err != nil {
		handleError(rw, err)
		return
	}
	// Remove hashed-password from response.
	user.Password = ""
	Created(rw, user)
}

func userChangePassword(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "User change password!")
}

func userLogin(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "User Login!")
}

func userLogout(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "User Logout!")
}
