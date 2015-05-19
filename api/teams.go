package api

import (
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

type TeamsHandler struct {
	ApiHandler
}

func (handler *TeamsHandler) CreateTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	team := &Team{}
	err = handler.parseBody(r.Body, team)
	if err != nil {
		return handler.handleError(err)
	}
	err = team.Save(currentUser)
	if err != nil {
		return handler.handleError(err)
	}
	return Created(team.ToString())
}

func (handler *TeamsHandler) UpdateTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	err = handler.parseBody(r.Body, team)
	if err != nil {
		return handler.handleError(err)
	}
	err = team.Save(currentUser)
	if err != nil {
		return handler.handleError(err)
	}
	return OK(team.ToString())
}

func (handler *TeamsHandler) DeleteTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := DeleteTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		return handler.handleError(err)
	}
	return OK(team.ToString())
}

func (handler *TeamsHandler) GetUserTeams(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	teams, _ := currentUser.GetTeams()
	s := CollectionSerializer{Items: teams, Count: len(teams)}
	payload := s.Serializer()
	return OK(payload)
}

func (handler *TeamsHandler) GetTeamInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	return OK(team.ToString())
}

func (handler *TeamsHandler) AddUsersToTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	var t *Team
	err = handler.parseBody(r.Body, &t)
	if err != nil {
		return handler.handleError(err)
	}
	err = team.AddUsers(t.Users)
	if err != nil {
		return handler.handleError(err)
	}
	return OK(team.ToString())
}

func (handler *TeamsHandler) RemoveUsersFromTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	var t *Team
	err = handler.parseBody(r.Body, &t)
	if err != nil {
		return handler.handleError(err)
	}
	err = team.RemoveUsers(t.Users)
	if err != nil {
		return Forbidden(err.Error())
	}
	return OK(team.ToString())
}
