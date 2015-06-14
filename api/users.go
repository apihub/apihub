package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
)

func (api *Api) userSignup(rw http.ResponseWriter, r *http.Request) {
	user := account.User{}
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

func (api *Api) userDelete(rw http.ResponseWriter, r *http.Request) {
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

func (api *Api) userChangePassword(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		account.User
		NewPassword          string `json:"new_password,omitempty"`
		ConfirmationPassword string `json:"confirmation_password,omitempty"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	if u.NewPassword != u.ConfirmationPassword || u.NewPassword == "" {
		handleError(rw, errors.ErrConfirmationPassword)
		return
	}
	authUser, ok := api.auth.Authenticate(u.Email, u.Password)
	if !ok {
		handleError(rw, errors.ErrAuthenticationFailed)
		return
	}

	authUser.Password = u.NewPassword
	if err := authUser.ChangePassword(); err != nil {
		handleError(rw, err)
		return
	}

	NoContent(rw)
}

func (api *Api) userLogin(rw http.ResponseWriter, r *http.Request) {
	user := account.User{}
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

func (api *Api) userLogout(rw http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	if authToken != "" {
		api.auth.RevokeUserToken(authToken)
	}

	NoContent(rw)
}

// Split Authenticate and CreateUserToken because we can override only the authentication method and still use the token method.
func (api *Api) Login(email, password string) (*account.TokenInfo, error) {
	user, ok := api.auth.Authenticate(email, password)
	if ok {
		token, err := api.auth.CreateUserToken(user)
		if err != nil {
			Logger.Warn(err.Error())
			return nil, err
		}
		return token, nil
	}

	return nil, errors.ErrAuthenticationFailed
}
