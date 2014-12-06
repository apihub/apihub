package account

import (
	"strings"

	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Service struct {
	Subdomain       string `bson:"_id" json:"subdomain"`
	AllowKeylessUse bool   `json:"allow_keyless_use"`
	Description     string `json:"description"`
	Disabled        bool   `json:"disabled"`
	Documentation   string `json:"documentation"`
	Endpoint        string `json:"endpoint"`
	Owner           string `json:"owner"`
	Team            string `json:"team"`
	Timeout         int    `json:"timeout"`
}

func (service *Service) Save(owner *User, team *Team) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if service.Subdomain == "" {
		message := "Subdomain cannot be empty."
		return &errors.ValidationError{Message: message}
	}
	if len(service.Endpoint) == 0 {
		message := "Endpoint cannot be empty."
		return &errors.ValidationError{Message: message}
	}

	service.Subdomain = strings.ToLower(service.Subdomain)
	service.Owner = owner.Email
	service.Team = team.Alias

	err = conn.Services().Insert(service)
	if mgo.IsDup(err) {
		message := "There is another service with this subdomain."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func (service *Service) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Services().Remove(service)
	if err == mgo.ErrNotFound {
		message := "Document not found."
		return &errors.ValidationError{Message: message}
	}
	if err != nil {
		return &errors.ValidationError{Message: err.Error()}
	}
	return err
}

func CountService() (int, error) {
	conn, err := db.Conn()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return conn.Services().Count()
}

func FindServiceBySubdomain(subdomain string) (*Service, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var service Service
	err = conn.Services().FindId(subdomain).One(&service)
	if err == mgo.ErrNotFound {
		message := "Service not found."
		return nil, &errors.ValidationError{Message: message}
	}

	return &service, nil
}

func DeleteServiceBySubdomain(subdomain string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Services().Remove(bson.M{"_id": subdomain})
	if err == mgo.ErrNotFound {
		message := "Service not found."
		return &errors.ValidationError{Message: message}
	}

	return nil
}
