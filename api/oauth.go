package api

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"runtime"

	"github.com/RangelReale/osin"
	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
	. "github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
	"github.com/fatih/structs"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

type OAuthHandler struct {
	Handler
}

type PageForm struct {
	Action             string
	Client             *Client
	InvalidCredentials bool
	Data               map[string]interface{}
	Method             string
}

type AuthenticationInfo struct {
	ClientId     string `json:"client_id"`
	Expires      int    `json:"expires_in"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"token_type"`
	User         string `json:"user"`
}

func (a *AuthenticationInfo) save() {
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	tokenKey := a.Type + " " + a.Token
	go conn.Tokens(tokenKey, a.Expires, structs.Map(a))
}

func (handler *OAuthHandler) Token(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	api := c.Env["Api"].(*Api)
	resp := api.oAuthServer.NewResponse()
	defer resp.Close()

	ar := api.oAuthServer.HandleAccessRequest(resp, r)
	if ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
			Logger.Debug("Grant Type: Authorization Code")
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
			Logger.Debug("Grant Type: Refresh Token")
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
			Logger.Debug("Grant Type: Client Credentials")
		}
		api.oAuthServer.FinishAccessRequest(resp, r, ar)
		saveAuthenticationAccessInfo(ar, resp)
	}
	if resp.IsError && resp.InternalError != nil {
		Logger.Error("ERROR: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
	return nil
}

func saveAuthenticationAccessInfo(ar *osin.AccessRequest, resp *osin.Response) {
	go func() {
		auth := &AuthenticationInfo{}
		auth.ClientId = ar.Client.GetId()
		// There is no UserData for client credentials.
		if ar.UserData != nil {
			u := ar.UserData.(bson.M)
			auth.User = u["email"].(string)
		}
		auth.Token = resp.Output["access_token"].(string)
		auth.Expires = int(resp.Output["expires_in"].(int32))
		auth.Type = resp.Output["token_type"].(string)
		// There is no refresh_token for client credentials.
		if resp.Output["refresh_token"] != nil {
			auth.RefreshToken = resp.Output["refresh_token"].(string)
		}
		auth.save()
	}()
}

func HandleLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) *User {
	r.ParseForm()
	p := &PageForm{
		Action: fmt.Sprintf("/login/oauth/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s", ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri)),
		Method: "POST",
	}

	if client, err := FindClientById(ar.Client.GetId()); err == nil {
		p.Data = map[string]interface{}{"client": client.Name}
	}

	if r.Method == "POST" {
		user := &User{Email: r.Form.Get("email"), Password: r.Form.Get("password")}
		if u, err := Login(user); err == nil {
			return u
		}
		p.InvalidCredentials = true
	}

	_, filename, _, _ := runtime.Caller(1)
	dir := path.Join(path.Dir(filename), "views/login.html")
	w.Header().Set("Content-Type", "text/html")
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
		Logger.Warn("Could not find client: %s.", resp.InternalError)
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
