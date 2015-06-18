package api

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	"github.com/gorilla/mux"
)

func serviceCreate(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	service := account.Service{}
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	team, err := findTeamAndCheckUser(service.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := service.Create(*user, *team); err != nil {
		handleError(rw, err)
		return
	}

	Created(rw, service)
}

func serviceUpdate(rw http.ResponseWriter, r *http.Request) {
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

	_, err = findTeamAndCheckUser(service.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}
	// It is not allowed to change the subdomain yet.
	service.Subdomain = mux.Vars(r)["subdomain"]

	err = service.Update()
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, service)
}

func serviceDelete(rw http.ResponseWriter, r *http.Request) {
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

	if err = service.Delete(*user); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, service)
}

func serviceInfo(rw http.ResponseWriter, r *http.Request) {
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

	_, err = findTeamAndCheckUser(service.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, service)
}

func serviceList(rw http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		handleError(rw, err)
		return
	}

	services, _ := user.Services()
	Ok(rw, CollectionSerializer{Items: services, Count: len(services)})
}
