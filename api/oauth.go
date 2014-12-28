package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/RangelReale/osin"
	. "github.com/backstage/backstage/account"
	"github.com/zenazn/goji/web"
)

type OAuthHandler struct {
	ApiHandler
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

	if r.Method == "POST" {
		user := &User{Email: r.Form.Get("email"), Password: r.Form.Get("password")}
		if u, err := Login(user); err == nil {
			return u
		}
	}

	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s\" method=\"POST\">",
		ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri))))

	w.Write([]byte("Email: <input type=\"text\" name=\"email\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

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
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
	return nil
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
