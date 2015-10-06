package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	. "github.com/apihub/apihub/log"
	"github.com/gorilla/mux"
)

func (api *Api) serviceCreate(rw http.ResponseWriter, r *http.Request, user *account.User) {
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

	go api.EventNotifier(newServiceEvent("service.create", service))
	Created(rw, service)
}

func (api *Api) serviceUpdate(rw http.ResponseWriter, r *http.Request, user *account.User) {
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

	go api.EventNotifier(newServiceEvent("service.update", *service))
	Ok(rw, service)
}

func (api *Api) serviceDelete(rw http.ResponseWriter, r *http.Request, user *account.User) {
	service, err := account.FindServiceBySubdomain(mux.Vars(r)["subdomain"])
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = service.Delete(*user); err != nil {
		handleError(rw, err)
		return
	}

	go api.EventNotifier(newServiceEvent("service.delete", *service))
	Ok(rw, service)
}

func (api *Api) serviceInfo(rw http.ResponseWriter, r *http.Request, user *account.User) {
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

func (api *Api) serviceList(rw http.ResponseWriter, r *http.Request, user *account.User) {
	services, _ := user.Services()
	Ok(rw, CollectionSerializer{Items: services, Count: len(services)})
}

type serviceEvent struct {
	CreatedAt time.Time       `json:"created_at"`
	Title     string          `json:"name"`
	Service   account.Service `json:"service"`
}

func (e *serviceEvent) Name() string {
	return e.Title
}

func (e *serviceEvent) Data() []byte {
	j, err := json.Marshal(e)
	if err != nil {
		Logger.Error(fmt.Sprintf("Failed to create a service event: %+v.", err))
	}
	return j
}

func newServiceEvent(name string, service account.Service) *serviceEvent {
	return &serviceEvent{
		CreatedAt: time.Now().UTC(),
		Title:     name,
		Service:   service,
	}
}
