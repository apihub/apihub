package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	"github.com/gorilla/mux"
)

func (api *Api) teamCreate(rw http.ResponseWriter, r *http.Request, user *account.User) {
	team := account.Team{}
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	if err := team.Create(*user); err != nil {
		handleError(rw, err)
		return
	}

	Created(rw, team)
}

func (api *Api) teamUpdate(rw http.ResponseWriter, r *http.Request, user *account.User) {
	team, err := findTeamAndCheckUser(mux.Vars(r)["alias"], user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	if err := team.Update(); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, team)
}
func (api *Api) teamList(rw http.ResponseWriter, r *http.Request, user *account.User) {
	teams, _ := user.Teams()
	Ok(rw, CollectionSerializer{Items: teams, Count: len(teams)})
}

func (api *Api) teamDelete(rw http.ResponseWriter, r *http.Request, user *account.User) {
	team, err := account.FindTeamByAlias(mux.Vars(r)["alias"])
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = team.Delete(*user); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, team)
}

func (api *Api) teamInfo(rw http.ResponseWriter, r *http.Request, user *account.User) {
	team, err := findTeamAndCheckUser(mux.Vars(r)["alias"], user)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, team)
}

func (api *Api) teamAddUsers(rw http.ResponseWriter, r *http.Request, user *account.User) {
	team, err := findTeamAndCheckUser(mux.Vars(r)["alias"], user)
	if err != nil {
		handleError(rw, err)
		return
	}

	var t *account.Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	err = team.AddUsers(t.Users)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, team)
}

func (api *Api) teamRemoveUsers(rw http.ResponseWriter, r *http.Request, user *account.User) {
	team, err := findTeamAndCheckUser(mux.Vars(r)["alias"], user)
	if err != nil {
		handleError(rw, err)
		return
	}

	var t *account.Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}
	err = team.RemoveUsers(t.Users)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, team)
}

func findTeamAndCheckUser(alias string, user *account.User) (*account.Team, error) {
	team, err := account.FindTeamByAlias(alias)
	if err != nil {
		return nil, err
	}
	if _, err := team.ContainsUser(user); err != nil {
		return nil, err
	}

	return team, nil
}
