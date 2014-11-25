package app

import (
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
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
	defer conn.Close()

	if user.Name == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		message := "Name/Email/Username/Password cannot be empty."
		return &errors.ValidationError{Message: message}
	}

	user.HashPassword()
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
	defer conn.Close()

	err = conn.Users().Remove(user)
	if err != nil && strings.Contains(err.Error(), "not found") {
		message := "User not found."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func FindUserByUsername(username string) (*User, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var result User
	err = conn.Users().Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return nil, &errors.ValidationError{Message: "User not found"}
	}

	return &result, nil
}

func (user *User) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user.Password = string(hash[:])
}