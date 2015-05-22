package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

type PluginsHandler struct {
	ApiHandler
}

func (handler *PluginsHandler) SubscribePlugin(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	conf := &PluginConfig{
		Name: c.URLParams["name"],
	}
	err = handler.parseBody(r.Body, conf)
	if err != nil {
		return handler.handleError(err)
	}

	err = conf.Save(currentUser)
	if err != nil {
		return handler.handleError(err)
	}

	payload, _ := json.Marshal(conf)
	return OK(string(payload))
}

func (handler *PluginsHandler) UnsubscribePlugin(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	currentUser, err := handler.getCurrentUser(c)
	if err != nil {
		return handler.handleError(err)
	}

	conf := &PluginConfig{
		Name: c.URLParams["name"],
	}
	err = handler.parseBody(r.Body, conf)
	if err != nil {
		return handler.handleError(err)
	}

	err = conf.Delete(currentUser)
	if err != nil {
		return handler.handleError(err)
	}
	payload, _ := json.Marshal(conf)
	return OK(string(payload))
}
