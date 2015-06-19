package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	"github.com/gorilla/mux"
)

func pluginSubsribe(rw http.ResponseWriter, r *http.Request, user *account.User) {
	service, err := account.FindServiceBySubdomain(mux.Vars(r)["subdomain"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = findTeamAndCheckUser(service.Team, user)
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

func pluginUnsubsribe(rw http.ResponseWriter, r *http.Request, user *account.User) {
	service, err := account.FindServiceBySubdomain(mux.Vars(r)["subdomain"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = findTeamAndCheckUser(service.Team, user)
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
