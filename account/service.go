// Package account encapsulates the business logic, determing how to manage teams, users and services.
package account

import (
	"encoding/json"
	"strings"

	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const CHANNEL_NAME = "services"

// The Service type is an encapsulation of a service details.
// It is not allowed to have more than one service with the same subdomain.
// The `AllowKeylessUse` field indicates if the proxy should validate the authorization header.
// The `Disabled` field indicates if the proxy dispatch the requests to the service.
// The `Timeout` field represents how milliseconds the proxy wait for the response, before returning an error.
type Service struct {
	Subdomain     string   `bson:"_id" json:"subdomain"`
	Description   string   `json:"description"`
	Disabled      bool     `json:"disabled"`
	Documentation string   `json:"documentation"`
	Endpoint      string   `json:"endpoint"`
	Transformers  []string `json:"transformers,omitempty"`
	Owner         string   `json:"owner"`
	Team          string   `json:"team"`
	Timeout       int      `json:"timeout"`
}

// Save creates a new service.
//
// It requires to inform the fields: Subdomain and Endpoint.
// It is not allowed to create two services with the same subdomain.
func (service *Service) Save(owner *User, team *Team) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if service.Subdomain == "" {
		return &errors.ValidationError{Payload: "Subdomain cannot be empty."}
	}
	if service.Endpoint == "" {
		return &errors.ValidationError{Payload: "Endpoint cannot be empty."}
	}

	service.Subdomain = strings.ToLower(service.Subdomain)
	service.Owner = owner.Email
	service.Team = team.Alias

	es, err := FindServiceBySubdomain(service.Subdomain)
	if err == nil && service.Team == es.Team {
		_, err = conn.Services().UpsertId(service.Subdomain, service)
	} else {
		err = conn.Services().Insert(service)
	}
	if mgo.IsDup(err) {
		return &errors.ValidationError{Payload: "There is another service with this subdomain."}
	}
	go service.publish()
	return err
}

// Delete removes an existing service from the server.
func (service *Service) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Services().Remove(service)
	if err == mgo.ErrNotFound {
		return &errors.ValidationError{Payload: "Service not found."}
	}
	if err != nil {
		return &errors.ValidationError{Payload: err.Error()}
	}
	go service.unpublish()
	return err
}

// Return the total of services in the database.
func CountService() (int, error) {
	conn, err := db.Conn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return conn.Services().Count()
}

// Try to find a service by its subdomain.
// If the service is not found, return an error. Return the service otherwise.
func FindServiceBySubdomain(subdomain string) (*Service, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var service Service
	err = conn.Services().FindId(subdomain).One(&service)
	if err == mgo.ErrNotFound {
		return nil, &errors.NotFoundError{Payload: "Service not found."}
	}

	return &service, nil
}

// DeleteServicesBySubdomain removes an existing service from the server based on given subdomain.
func DeleteServiceBySubdomain(subdomain string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	service, err := FindServiceBySubdomain(subdomain)
	if err != nil {
		return &errors.ValidationError{Payload: "Service not found."}
	}
	go service.unpublish()
	err = conn.Services().Remove(bson.M{"_id": subdomain})
	if err == mgo.ErrNotFound {
		return &errors.ValidationError{Payload: "Service not found."}
	}

	return nil
}

func DeleteServicesByTeam(team string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		services, err := FindServicesByTeam([]string{team})
		if err != nil {
			return
		}
		for _, service := range services {
			service.unpublish()
		}
	}()

	_, err = conn.Services().RemoveAll(bson.M{"team": team})
	if err != nil {
		return err
	}

	return nil
}

// Find all the services for a given team alias.
// Return an empty list if nothing is found.
func FindServicesByTeam(teams []string) ([]*Service, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var services []*Service = []*Service{}
	err = conn.Services().Find(bson.M{"team": bson.M{"$in": teams}}).All(&services)
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (service *Service) Middlewares() ([]*MiddlewareConfig, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var midds []*MiddlewareConfig = []*MiddlewareConfig{}
	err = conn.MiddlewaresConfig().Find(bson.M{"service": service.Subdomain}).All(&midds)
	if err != nil {
		return nil, err
	}
	return midds, nil
}

func (service *Service) publish() {
	s, err := json.Marshal(service)
	if err != nil {
		panic(err)
	}
	go func() {
		cli := db.NewRedisClient()
		cli.Publish(CHANNEL_NAME, string(s))
		defer cli.Close()
	}()
}

func (service *Service) unpublish() {
	service.Disabled = true
	service.publish()
}
