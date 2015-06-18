package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	"github.com/gorilla/mux"
)

func webhookSave(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	webhook := account.Webhook{}
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	team, err := findTeamAndCheckUser(webhook.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := webhook.Save(*team); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, webhook)
}

func webhookDelete(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	webhook, err := account.FindWebhookByName(mux.Vars(r)["name"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = findTeamAndCheckUser(webhook.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = webhook.Delete(); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, webhook)
}
