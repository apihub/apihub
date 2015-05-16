package account

import (
	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MiddlewareConfig struct {
	Name    string                 `json:"name"`
	Service string                 `json:"-""`
	Config  map[string]interface{} `json:"config"`
}

// Save associates a middleware with a service.
//
// It requires to inform the fields: Subdomain and Middleware name.
// It is not allowed to associates two middlewares with the same subdomain.
func (m *MiddlewareConfig) Save() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if m.Name == "" {
		return &errors.ValidationError{Payload: "Name cannot be empty."}
	}
	if m.Service == "" {
		return &errors.ValidationError{Payload: "Service cannot be empty."}
	}

	_, err = conn.MiddlewaresConfig().Upsert(bson.M{"service": m.Service, "name": m.Name}, m)
	return err
}

// Delete removes an existing middleware config from the server.
func (m *MiddlewareConfig) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.MiddlewaresConfig().Remove(m)
	if err == mgo.ErrNotFound {
		return &errors.ValidationError{Payload: "Middleware Config not found."}
	}
	if err != nil {
		return &errors.ValidationError{Payload: err.Error()}
	}
	return err
}
