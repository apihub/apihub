package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/backstage/backstage/db"
	"github.com/zenazn/goji/web"
	. "gopkg.in/check.v1"
)

func (s *S) TestTokenForAuthorizationCode(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)

	defer s.oAuthStorage.RemoveAuthorize(authorizeData.Code)
	err := s.oAuthStorage.SaveAuthorize(authorizeData)
	c.Assert(err, IsNil)
	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	s.router.Post("/login/oauth/token", s.Api.route(oAuthHandler, "Token"))
	req, _ := http.NewRequest("POST", fmt.Sprintf("/login/oauth/token?grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s", osinClient.Id, url.QueryEscape(osinClient.RedirectUri), authorizeData.Code), strings.NewReader(""))
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Matches, "^{\"access_token\":\".*?\",\"expires_in\":3600,\"refresh_token\":\".*?\",\"token_type\":\"Bearer\"}$")
}

func (s *S) TestTokenForRefreshToken(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)

	defer s.oAuthStorage.RemoveAccess(accessData.AccessToken)
	err := s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	s.router.Post("/login/oauth/token", s.Api.route(oAuthHandler, "Token"))
	req, _ := http.NewRequest("POST", fmt.Sprintf("/login/oauth/token?grant_type=refresh_token&client_id=%s&redirect_uri=%s&refresh_token=%s", osinClient.Id, url.QueryEscape(osinClient.RedirectUri), accessData.RefreshToken), strings.NewReader(""))
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Matches, "^{\"access_token\":\".*?\",\"expires_in\":3600,\"refresh_token\":\".*?\",\"token_type\":\"Bearer\"}$")
}

func (s *S) TestTokenForClientCredentials(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)

	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	b := strings.NewReader("grant_type=client_credentials")
	s.router.Post("/login/oauth/token", s.Api.route(oAuthHandler, "Token"))
	req, _ := http.NewRequest("POST", "/login/oauth/token", b)
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Matches, "^{\"access_token\":\".*?\",\"expires_in\":.*?,\"token_type\":\"Bearer\"}$")
}

func (s *S) TestTokenForClientCredentialsWithoutBasicAuth(c *C) {
	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	b := strings.NewReader("grant_type=client_credentials")
	s.router.Post("/login/oauth/token", s.Api.route(oAuthHandler, "Token"))
	req, _ := http.NewRequest("POST", "/login/oauth/token", b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Equals, `{"error":"invalid_request","error_description":"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed."}`)
}

func (s *S) TestTokenForClientCredentialsWithInvalidCredentials(c *C) {
	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	b := strings.NewReader("grant_type=client_credentials")
	s.router.Post("/login/oauth/token", s.Api.route(oAuthHandler, "Token"))
	req, _ := http.NewRequest("POST", "/login/oauth/token", b)
	req.SetBasicAuth("invalid-id", "invalid-secret")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Equals, `{"error":"server_error","error_description":"The authorization server encountered an unexpected condition that prevented it from fulfilling the request."}`)
}

func (s *S) TestTokenForUnsupporteGrantType(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)

	defer s.oAuthStorage.RemoveAuthorize(authorizeData.Code)
	err := s.oAuthStorage.SaveAuthorize(authorizeData)
	c.Assert(err, IsNil)
	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	s.router.Post("/login/oauth/token", s.Api.route(oAuthHandler, "Token"))
	req, _ := http.NewRequest("POST", fmt.Sprintf("/login/oauth/token?grant_type=code&client_id=%s&redirect_uri=%s&code=%s", osinClient.Id, url.QueryEscape(osinClient.RedirectUri), authorizeData.Code), strings.NewReader(""))
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Equals, `{"error":"unsupported_grant_type","error_description":"The authorization grant type is not supported by the authorization server."}`)
}

func (s *S) TestInfoFromAccessTokenForClientCredentials(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)

	defer s.oAuthStorage.RemoveAccess(accessData.AccessToken)
	err := s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	s.router.Get("/me", s.Api.route(oAuthHandler, "Info"))
	req, _ := http.NewRequest("GET", fmt.Sprintf("/me?code=%s", accessData.AccessToken), nil)
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusOK)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Matches, `{"access_token":"test-123456","client_id":"test-1234","expires_in":.*?,"refresh_token":"test-refresh-7890","token_type":"Bearer"}`)
}

func (s *S) TestInfoFromAccessTokenForInvalidToken(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()

	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	s.router.Get("/me", s.Api.route(oAuthHandler, "Info"))
	req, _ := http.NewRequest("GET", "/me?code=invalid-token", nil)
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/json")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusBadRequest)
	c.Assert(strings.TrimRight(s.recorder.Body.String(), "\n"), Equals, `{"error":"invalid_request","error_description":"The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed."}`)
}

func (s *S) TestAuthorizeWithValidCredentials(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	alice.Save()
	defer alice.Delete()
	s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)

	s.Api.LoadOauthServer(s.oAuthStorage)
	s.env["Api"] = s.Api
	s.router.Post("/login/oauth/authorize", s.Api.route(oAuthHandler, "Authorize"))
	req, _ := http.NewRequest("POST", fmt.Sprintf("/login/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s", osinClient.Id, url.QueryEscape(osinClient.RedirectUri)), strings.NewReader(fmt.Sprintf("email=%s&password=123456", alice.Email)))
	req.SetBasicAuth(osinClient.Id, osinClient.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	webC := web.C{Env: s.env}
	s.router.ServeHTTPC(webC, s.recorder, req)

	c.Assert(s.recorder.Code, Equals, http.StatusFound)
}
