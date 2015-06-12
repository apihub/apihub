package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	"github.com/gorilla/mux"
)

func teamCreate(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

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

func teamList(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	teams, _ := user.Teams()
	Ok(rw, CollectionSerializer{Items: teams, Count: len(teams)})
}

func teamDelete(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

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

func teamInfo(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	team, err := account.FindTeamByAlias(mux.Vars(r)["alias"])
	if err != nil {
		handleError(rw, err)
		return
	}
	if _, err := team.ContainsUser(user); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, team)
}

func teamAddUsers(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	team, err := findTeamByAlias(mux.Vars(r)["alias"], user)
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

func teamRemoveUsers(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	team, err := findTeamByAlias(mux.Vars(r)["alias"], user)
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

func teamUpdate(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	team, err := findTeamByAlias(mux.Vars(r)["alias"], user)
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

func findTeamByAlias(alias string, user *account.User) (*account.Team, error) {
	team, err := account.FindTeamByAlias(alias)
	if err != nil {
		return nil, err
	}
	if _, err := team.ContainsUser(user); err != nil {
		return nil, err
	}

	return team, nil
}
