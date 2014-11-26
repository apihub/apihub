package app

import (
	"strings"
	"time"

	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
)

type Service struct {
	Subdomain       string `bson:"_id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	AllowKeylessUse bool
	Description     string
	Disabled        bool
	Documentation   string
	Endpoint        map[string]interface{}
	Owner           string
	Timeout         int
	Name            string
}

func CreateService(service *Service, user *User) error {
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

	// FIXME: improve this.
	service.Subdomain = strings.ToLower(service.Subdomain)
	service.CreatedAt = time.Now().In(time.UTC)
	service.UpdatedAt = time.Now().In(time.UTC)
	service.Owner = user.Username

	err = conn.Services().Insert(service)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		message := "There is another service with this subdomain."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func DeleteService(service *Service) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Services().Remove(service)
	if err != nil && strings.Contains(err.Error(), "not found") {
		message := "Document not found."
		return &errors.ValidationError{Message: message}
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
