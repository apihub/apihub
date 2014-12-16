package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type TeamsHandler struct {
	ApiHandler
}

func (handler *TeamsHandler) CreateTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	team := &Team{}
	err = handler.parseBody(r.Body, team)
	if err != nil {
		return BadRequest(err.Error())
	}
	err = team.Save(currentUser)
	if err != nil {
		return BadRequest(err.Error())
	}
	return Created(team.ToString())
}

func (handler *TeamsHandler) DeleteTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	team, err := DeleteTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		switch err.(type) {
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(err.Error())
		}
	}
	return OK(team.ToString())
}

func (handler *TeamsHandler) GetUserTeams(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	teams, _ := currentUser.GetTeams()
	payload, _ := json.Marshal(teams)
	return OK(string(payload))
}

func (handler *TeamsHandler) GetTeamInfo(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		switch err.(type) {
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(err.Error())
		}
	}

	return OK(team.ToString())
}

func (handler *TeamsHandler) AddUsersToTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		switch err.(type) {
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(err.Error())
		}
	}

	var t *Team
	err = handler.parseBody(r.Body, &t)
	if err != nil {
		return BadRequest(err.Error())
	}
	err = team.AddUsers(t.Users)
	if err != nil {
		return BadRequest(err.Error())
	}
	return Created(team.ToString())
}

func (handler *TeamsHandler) RemoveUsersFromTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	team, err := FindTeamByAlias(c.URLParams["alias"], currentUser)
	if err != nil {
		switch err.(type) {
		case *ForbiddenError:
			return Forbidden(err.Error())
		default:
			return BadRequest(err.Error())
		}
	}

	var t *Team
	err = handler.parseBody(r.Body, &t)
	if err != nil {
		return BadRequest(err.Error())
	}
	err = team.RemoveUsers(t.Users)
	if err != nil {
		return Forbidden(err.Error())
	}
	return OK(team.ToString())
}
