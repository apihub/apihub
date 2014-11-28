package system

import (
  "fmt"
  "encoding/json"
  "net/http"

  "github.com/zenazn/goji/web"
  "github.com/albertoleal/backstage/auth"
  "github.com/albertoleal/backstage/errors"
  "github.com/albertoleal/backstage/api/context"
)

func AuthorizationMiddleware(c *web.C, h http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    authorization := r.Header.Get("Authorization")
    _, _, err := auth.GetToken(authorization)
    if err != nil {
      context.AddRequestError(c, &errors.HTTPError{StatusCode: http.StatusUnauthorized, Message: "You do not have access to this resource."})
    }
    h.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func ErrorHandlerMiddleware(c *web.C, h http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    key, ok := context.GetRequestError(c)
    if ok {
      body, _ := json.Marshal(key)
      http.Error(w, string(body), key.StatusCode)
    } else {
      h.ServeHTTP(w, r)
    }
  }

  return http.HandlerFunc(fn)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  notFound := &errors.HTTPError{StatusCode: http.StatusNotFound, Message: "The resource you are looking for was not found."}
  w.WriteHeader(notFound.StatusCode)
  body, _ := json.Marshal(notFound)
  fmt.Fprint(w, string(body))
  return
}