package api

import (
	"encoding/json"
	"net/http"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	"github.com/gorilla/mux"
)

func (api *Api) hookSave(rw http.ResponseWriter, r *http.Request, user *account.User) {
	hook := account.Hook{}
	if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	team, err := findTeamAndCheckUser(hook.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := hook.Save(*team); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, hook)
}

func (api *Api) hookDelete(rw http.ResponseWriter, r *http.Request, user *account.User) {
	hook, err := account.FindHookByName(mux.Vars(r)["name"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = findTeamAndCheckUser(hook.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = hook.Delete(); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, hook)
}
