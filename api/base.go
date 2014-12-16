package api

import (
	"encoding/json"
	"io"
	"io/ioutil"

	. "github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

type ApiHandler struct{}

func (api *ApiHandler) getCurrentUser(c *web.C) (*User, error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (api *ApiHandler) parseBody(body io.ReadCloser, r interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &r); err != nil {
		return ErrBadRequest
	}
	return nil
}
