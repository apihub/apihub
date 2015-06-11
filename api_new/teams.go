package api_new

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
)

func teamCreate(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	team := account_new.Team{}
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
