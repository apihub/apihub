=======
Plugins
=======

It's possible to add functionalities to your APIs just by using plugins. Backstage supports two type: Middleware and Transformer.


Middleware
----------
Middleware is a wrapper around your API that decorates the requests without adding logic in the application. It's supposed to run before dispatching the request to the API. It's allowed to use as many middlewares as you want, since they implement the interface below:

.. image:: middleware.png
   :name: middleware

.. highlight:: go

::

  type Middleware interface {
    Configure(cfg string)
    Serve(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
  }



Writting a Middleware
~~~~~~~~~~~~~~~~~~~~~

.. highlight:: go

::

  package middleware

  import (
    "encoding/json"
    "net/http"

    "github.com/rs/cors"
  )

  type Cors struct {
    AllowedOrigins   []string `json:"allowed_origins"`
    AllowedMethods   []string `json:"allowed_methods"`
    AllowedHeaders   []string `json:"allowed_headers"`
    ExposedHeaders   []string `json:"exposed_headers"`
    AllowCredentials bool     `json:"allow_credentials"`
    MaxAge           int      `json:"max_age"`
    Debug            bool     `json:"debug"`
  }

  func NewCorsMiddleware() Middleware {
    return &Cors{}
  }

  func (c *Cors) Serve(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    cors := cors.New(cors.Options{
      AllowedOrigins:   c.AllowedOrigins,
      AllowedMethods:   c.AllowedMethods,
      AllowedHeaders:   c.AllowedHeaders,
      ExposedHeaders:   c.ExposedHeaders,
      AllowCredentials: c.AllowCredentials,
      MaxAge:           c.MaxAge,
      Debug:            c.Debug,
    })
    cors.ServeHTTP(rw, r, next)
  }

  func (c *Cors) Configure(cfg string) {
    json.Unmarshal([]byte(cfg), c)
  }


After that, it's needed to add the middleware to the Gateway:

.. highlight:: go

::

  gw.Middleware().Add("cors", NewCorsMiddleware)


The response:

.. highlight:: bash

::

  HTTP/1.1 200 OK
  access-control-allow-credentials: true
  access-control-allow-methods: PUT
  access-control-allow-origin: http://helloworld.backstage.dev
  access-control-max-age: 10
  vary: Origin
  date: Sat, 16 May 2015 13:40:44 GMT
  content-length: 0
  content-type: text/plain; charset=utf-8
  Connection: keep-alive

Using a Middleware
~~~~~~~~~~~~~~~~~~~~

To use a Middleware, it's needed to create a config for each service and it's needed to use the name you used when adding it to the Gateway:

.. highlight:: bash

::

  curl -XOPTIONS -H 'Access-Control-Request-Method: PUT' -H 'Origin: http://helloworld.backstage.dev' http://helloworld.backstage.dev/ -i

.. highlight:: go

::

  services := []*account.Service{&account.Service{Endpoint: "http://www.example.org", Subdomain: "example"}}
  confCors := &account.MiddlewareConfig{
    Name:    "cors",
    Service: services[0].Subdomain,
    Config:  map[string]interface{}{"allowed_origins": []string{"http://helloworld.backstage.dev"}, "debug": true, "allowed_methods": []string{"DELETE", "PUT"}, "allow_credentials": true, "max_age": 10},
  }
  confCors.Save()


Transformer
-----------
Transformer is supposed to run after the API response, just before writing the final response.

.. image:: transformer.png
   :name: transformer

.. highlight:: go

::

  type Filter func(*http.Request, *http.Response, *bytes.Buffer)

Writting a Transform
~~~~~~~~~~~~~~~~~~~~

.. highlight:: go

::

  func FooTransformer(r *http.Request, w *http.Response, body *bytes.Buffer) {
    w.Header.Set("Content-Type", "text/plain")
    body.Reset()
    body.Write([]byte("Foo"))
  }

After that, it's needed to add the transformer to the Gateway:

.. highlight:: go

::

  gateway.Transformer().Add("FooTransformer", FooTransformer)


Using a Transform
~~~~~~~~~~~~~~~~~~~~

To use a Transformer, you just need to use the name you used when adding it to the Gateway:

.. highlight:: go

::

  services := []*account.Service{&account.Service{Endpoint: "http://www.example.org", Subdomain: "example",Transformers: []string{"FooTransformer"}}}