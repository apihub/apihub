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

func userDelete(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := user.Delete(); err != nil {
		handleError(rw, err)
		return
	}
	// Remove hashed-password from response.
	user.Password = ""
	Ok(rw, user)
}

func userChangePassword(rw http.ResponseWriter, r *http.Request) {
	// NewPassword          string `json:"new_password,omitempty"`
	// ConfirmationPassword string `json:"confirmation_password,omitempty"`
	fmt.Fprintln(rw, "User change password!")
}

func (api *Api) userLogin(rw http.ResponseWriter, r *http.Request) {
	user := account_new.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	token, err := api.Login(user.Email, user.Password)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, token)
}

func userLogout(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "User Logout!")
}
