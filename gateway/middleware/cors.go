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

func (c *Cors) ProcessRequest(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
