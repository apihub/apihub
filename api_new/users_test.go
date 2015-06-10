package api_new_test

import (
	"fmt"
	"net/http"

	"github.com/backstage/backstage/auth_new"
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
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/api/users",
		Headers: http.Header{"Authorization": {"expired-token"}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
}

func (s *S) TestDeleteUserWithoutToken(c *C) {
	headers, code, body, err := httpClient.MakeRequest(RequestArgs{
		Method: "DELETE",
		Path:   "/api/users",
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(headers.Get("Content-Type"), Equals, "application/json")
	c.Assert(string(body), Equals, `{"error":"bad_request","error_description":"Invalid or expired token. Please log in with your Backstage credentials."}`)
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
	c.Assert(string(body), Matches, fmt.Sprintf(`{"access_token":".*","token_type":"%s","expires":%d,"created_at":".*"}`, auth_new.TOKEN_TYPE, auth_new.EXPIRES_IN_SECONDS))
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
	_, code, _, err := httpClient.MakeRequest(RequestArgs{
		Method:  "DELETE",
		Path:    "/auth/logout",
		Headers: http.Header{"Authorization": {s.authHeader}},
	})

	c.Check(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
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

}
