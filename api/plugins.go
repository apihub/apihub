package api

import (
	"encoding/json"
	"net/http"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	"github.com/gorilla/mux"
)

func (api *Api) pluginSubsribe(rw http.ResponseWriter, r *http.Request, user *account.User) {
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

	plugin := account.Plugin{}
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

func (api *Api) pluginUnsubsribe(rw http.ResponseWriter, r *http.Request, user *account.User) {
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
