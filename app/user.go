package app

import (
	"strings"

	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
)

type User struct {
	Name     string
	Email    string
	Username string
	Password string //[]byte
}

func CreateUser(user *User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}

	if user.Name == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		message := "Name/Email/Username/Password cannot be empty."
		return &errors.ValidationError{Message: message}
	}

	err = conn.Users().Insert(user)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		message := "Someone already has that username. Could you try another?."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func DeleteUser(user *User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}

	err = conn.Users().Remove(user)
	if err != nil && strings.Contains(err.Error(), "not found") {
		message := "User not found."
		return &errors.ValidationError{Message: message}
	}
	return err
}