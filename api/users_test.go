package api_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/backstage/auth"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateUser(c *C) {
	defer func() {
		store, _ := s.store()
		u, _ := store.FindUserByEmail("alice@example.org")
		u.Delete()
	}()

	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "POST",
		Path:   "/auth/signup",
		Body:   `{"name": "Alice", "email": "alice@example.org", "password": "123456"}`,
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"name":"Alice","email":"alice@example.org"}`)
}

func (s *S) TestCreateUserWithInvalidPayloadFormat(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "POST",
		Path:   "/auth/signup",
		Body:   `"name": "Alice"`,
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"The request was invalid or cannot be served."}`)
}

func (s *S) TestCreateUserWithMissingRequiredFields(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "POST",
		Path:   "/auth/signup",
		Body:   `{}`,
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Name/Email/Password cannot be empty."}`)
}

func (s *S) TestDeleteUser(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/users",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, fmt.Sprintf(`{"name":"%s","email":"%s"}`, user.Name, user.Email))
}

func (s *S) TestDeleteUserWithExpiredToken(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/users", Headers: http.Header{"Authorization": {"expired-token"}}}, c)
}

func (s *S) TestDeleteUserWithoutToken(c *C) {
	testWithoutSignIn(RequestArgs{Method: "DELETE", Path: "/api/users"}, c)
}

func (s *S) TestLoginUser(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "POST",
		Path:   "/auth/login",
		Body:   fmt.Sprintf(`{"email": "%s", "password": "secret"}`, user.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Matches, fmt.Sprintf(`{"access_token":".*","created_at":".*","expires":%d,"token_type":"%s"}`, auth.EXPIRES_IN_SECONDS, auth.TOKEN_TYPE))
}

func (s *S) TestLoginUserWithInvalidUser(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "POST",
		Path:   "/auth/login",
		Body:   fmt.Sprintf(`{"email": "%s", "password": "invalid-password"}`, user.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Authentication failed."}`)
}

func (s *S) TestLogoutUser(c *C) {
	_, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/auth/logout",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
	c.Assert(string(body), Equals, "")
}

func (s *S) TestLogoutUserWithInvalidToken(c *C) {
	_, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/auth/logout",
		Headers: http.Header{"Authorization": {"invalid-token"}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
	c.Assert(string(body), Equals, "")
}

func (s *S) TestChangePassword(c *C) {
	defer func() {
		store, _ := s.store()
		u, _ := store.FindUserByEmail("bob@bar.example.org")
		u.Delete()
	}()

	_, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "PUT",
		Path:   "/auth/password",
		Body:   fmt.Sprintf(`{"email": "%s", "password": "secret", "new_password": "123", "confirmation_password": "123"}`, user.Email),
	})

	c.Check(err, IsNil)
	c.Assert(string(body), Equals, "")
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *S) TestChangePasswordWithInvalidCredentials(c *C) {
	_, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "PUT",
		Path:   "/auth/password",
		Body:   fmt.Sprintf(`{"email": "%s", "password": "%s", "new_password": "123", "confirmation_password": "123"}`, user.Email, "invalid-password"),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Authentication failed."}`)
}

func (s *S) TestChangePasswordWithInvalidNewPassword(c *C) {
	_, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "PUT",
		Path:   "/auth/password",
		Body:   fmt.Sprintf(`{"email": "%s", "password": "secret"}`, user.Email),
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Your new password and confirmation password do not match or are invalid."}`)
}
