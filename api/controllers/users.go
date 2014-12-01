package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/api/context"
	"github.com/albertoleal/backstage/errors"

	"github.com/zenazn/goji/web"
)

type UsersController struct {
	ApiController
}

func (controller *UsersController) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) (string, int) {
	user := &User{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "Deu ruim!"})
	}
	if err = json.Unmarshal(body, user); err != nil {
		context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "Deu ruim!"})
	}
	ok := user.Save()
	if ok != nil {
		context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusBadRequest, Message: "Deu ruim!"})
	}
	user.Password = ""
	body, _ = json.Marshal(user)
	return string(body), http.StatusCreated
}
