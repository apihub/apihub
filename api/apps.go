package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	"github.com/gorilla/mux"
)

func appCreate(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	app := account.App{}
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	team, err := userBelongsToTeam(app.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := app.Create(*user, *team); err != nil {
		handleError(rw, err)
		return
	}

	Created(rw, app)
}

func appUpdate(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	app, err := account.FindAppByClientId(mux.Vars(r)["client_id"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = userBelongsToTeam(app.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}
	// It is not allowed to change the client id yet.
	app.ClientId = mux.Vars(r)["client_id"]

	err = app.Update()
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, app)
}

func appDelete(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	app, err := account.FindAppByClientId(mux.Vars(r)["client_id"])
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = app.Delete(*user); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, app)
}

func appInfo(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	app, err := account.FindAppByClientId(mux.Vars(r)["client_id"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = userBelongsToTeam(app.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, app)
}
