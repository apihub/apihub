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
	// TODO: Remove this logic from here.
	team, err := FindTeamByAlias(c.URLParams["alias"])
	if err != nil || team.Owner != currentUser.Email {
		return Forbidden(ErrOnlyOwnerHasPermission.Error())
	}
	err = team.Delete()
	if err != nil {
		return BadRequest(err.Error())
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

	//TODO: Remove this login from here.
	team, err := FindTeamByAlias(c.URLParams["alias"])
	if err != nil {
		return BadRequest(err.Error())
	}
	_, err = team.ContainsUser(currentUser)
	if err != nil {
		return Forbidden(err.Error())
	}
	return OK(team.ToString())
}

func (handler *TeamsHandler) AddUsersToTeam(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return BadRequest(err.Error())
	}

	//TODO: Remove this login from here.
	team, err := FindTeamByAlias(c.URLParams["alias"])
	if err != nil {
		return BadRequest(ErrTeamNotFound.Error())
	}
	_, err = team.ContainsUser(currentUser)
	if err != nil {
		return Forbidden(err.Error())
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

	//TODO: Remove this login from here
	team, err := FindTeamByAlias(c.URLParams["alias"])
	if err != nil {
		return BadRequest(ErrTeamNotFound.Error())
	}
	_, err = team.ContainsUser(currentUser)
	if err != nil {
		return Forbidden(err.Error())
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