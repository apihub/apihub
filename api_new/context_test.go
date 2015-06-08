package api_new_test

// import (
//   "encoding/json"
//   "fmt"
//   "net/http"
//   "net/http/httptest"

//   "github.com/backstage/backstage/account"
//   "github.com/backstage/backstage/errors"
//   "github.com/zenazn/goji/web"
//   . "gopkg.in/check.v1"
// )

// func (s *S) TestAddGetRequestError(c *C) {
//   m := web.New()

//   m.Get("/helloworld", func(c web.C, w http.ResponseWriter, r *http.Request) {
//     AddRequestError(&c, &HTTPResponse{StatusCode: http.StatusUnauthorized,
//       ErrorType:        errors.E_UNAUTHORIZED_REQUEST,
//       ErrorDescription: "You do not have access to this resource."})

//     key, _ := GetRequestError(&c)
//     body, _ := json.Marshal(key)
//     http.Error(w, string(body), key.StatusCode)
//   })

//   req, _ := http.NewRequest("GET", "/helloworld", nil)
//   recorder := httptest.NewRecorder()
//   env := map[string]interface{}{}
//   m.ServeHTTPC(web.C{Env: env}, recorder, req)

//   c.Assert(recorder.Code, Equals, 401)
//   c.Assert(recorder.Body.String(), Equals, "{\"error\":\"unauthorized_access\",\"error_description\":\"You do not have access to this resource.\"}\n")
// }

// func (s *S) TestSetAndGetCurrentUser(c *C) {
//   m := web.New()

//   m.Get("/helloworld", func(c web.C, w http.ResponseWriter, r *http.Request) {
//     alice := &account.User{Username: "alice", Name: "Alice", Email: "alice@example.org", Password: "123456"}
//     alice.Save()
//     defer alice.Delete()
//     SetCurrentUser(&c, alice)
//     user, _ := GetCurrentUser(&c)
//     body, _ := json.Marshal(user)
//     fmt.Fprint(w, string(body))
//     w.WriteHeader(http.StatusOK)
//   })

//   req, _ := http.NewRequest("GET", "/helloworld", nil)
//   recorder := httptest.NewRecorder()
//   env := map[string]interface{}{}
//   m.ServeHTTPC(web.C{Env: env}, recorder, req)

//   c.Assert(recorder.Code, Equals, http.StatusOK)
// }

// func (s *S) TestGetCurrentUserWhenNotSignedIn(c *C) {
//   m := web.New()

//   m.Get("/helloworld", func(co web.C, w http.ResponseWriter, r *http.Request) {
//     _, err := GetCurrentUser(&co)
//     c.Assert(err.Error(), Equals, "Invalid or expired token. Please log in with your Backstage credentials.")
//   })

//   req, _ := http.NewRequest("GET", "/helloworld", nil)
//   recorder := httptest.NewRecorder()
//   env := map[string]interface{}{}
//   m.ServeHTTPC(web.C{Env: env}, recorder, req)
// }
