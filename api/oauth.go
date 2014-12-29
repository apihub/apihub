package api

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/RangelReale/osin"
	. "github.com/backstage/backstage/account"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type OAuthHandler struct {
	ApiHandler
}

type PageForm struct {
	Action string
	Client *Client
	InvalidCredentials bool
	Data map[string]interface{}
	Method string
}

func (handler *OAuthHandler) Token(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	api := c.Env["Api"].(*Api)
	resp := api.oAuthServer.NewResponse()
	defer resp.Close()

	if ar := api.oAuthServer.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		}
		api.oAuthServer.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
	return nil
}

func HandleLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) *User {
	r.ParseForm()
	p := &PageForm{
		Action: fmt.Sprintf("/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s",ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri)),
		Method: "POST",
		Data: map[string]interface{}{"client": "tesssst"},
	}
	if client, err := FindClientById(ar.Client.GetId()); err == nil {
		p.Client = client
	}

	if r.Method == "POST" {
		user := &User{Email: r.Form.Get("email"), Password: r.Form.Get("password")}
		if u, err := Login(user); err == nil {
			return u
		}
		p.InvalidCredentials = true
	}

	dir, err := filepath.Abs("api/views/login.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t, _ := template.ParseFiles(dir)
  t.Execute(w, p)
	return nil
}

func (handler *OAuthHandler) Authorize(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	api := c.Env["Api"].(*Api)
	c.Env["Content-Type"] = "text/html"
	resp := api.oAuthServer.NewResponse()
	defer resp.Close()

	if ar := api.oAuthServer.HandleAuthorizeRequest(resp, r); ar != nil {
		user := HandleLoginPage(ar, w, r)
		if user == nil {
			return nil
		}
		ar.UserData = struct{ Username, Email, Name string }{Username: user.Username, Email: user.Email, Name: user.Name}
		ar.Authorized = true
		api.oAuthServer.FinishAuthorizeRequest(resp, r, ar)
	}
	//That's a hack to avoid redirection when redirect_uri does not match. Have opened an issue: RangelReale/osin/issues/41.
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
		return BadRequest(E_BAD_REQUEST, resp.Output["error_description"].(string))
	} else {
		osin.OutputJSON(resp, w, r)
		return nil
	}
}

func (handler *OAuthHandler) Info(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	api := c.Env["Api"].(*Api)
	resp := api.oAuthServer.NewResponse()
	defer resp.Close()
	ir := api.oAuthServer.HandleInfoRequest(resp, r)
	if ir != nil {
		api.oAuthServer.FinishInfoRequest(resp, r, ir)
	}
	if !resp.IsError {
		u := ir.AccessData.UserData
		if u != nil {
			resp.Output["user"] = u
		}
	}
	osin.OutputJSON(resp, w, r)
	return nil
}
