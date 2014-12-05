package api

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/zenazn/goji/web"
)

type ApiHandler struct{}

func (api *ApiHandler) getCurrentUser(c *web.C) (user *User, erro error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		erro := &HTTPResponse{StatusCode: http.StatusBadRequest, Payload: err.Error()}
		AddRequestError(c, erro)
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
		return errors.New("The request was bad-formed.")
	}
	return nil
}
