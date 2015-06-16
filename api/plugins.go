package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/apimanager/account"
	"github.com/backstage/apimanager/errors"
	"github.com/gorilla/mux"
)

func pluginSubsribe(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	service, err := account.FindServiceBySubdomain(mux.Vars(r)["subdomain"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = userBelongsToTeam(service.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	plugin := account.PluginConfig{}
	if err := json.NewDecoder(r.Body).Decode(&plugin); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	err = plugin.Save(*service)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, plugin)
}

func pluginUnsubsribe(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	service, err := account.FindServiceBySubdomain(mux.Vars(r)["subdomain"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = userBelongsToTeam(service.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	plugin, err := account.FindPluginByNameAndService(mux.Vars(r)["plugin_name"], *service)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = plugin.Delete(); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, plugin)
}
