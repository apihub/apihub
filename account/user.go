package account

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (user *User) Save() error {
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
	if mgo.IsDup(err) {
		message := "Someone already has that username. Could you try another?."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func (user *User) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Users().Remove(user)
	if err == mgo.ErrNotFound {
		message := "User not found."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func (user *User) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user.Password = string(hash[:])
}

func (user *User) Valid() bool {
	_, err := FindUserByUsername(user.Username)
	if err == nil {
		return true
	}
	return false
}

func FindUserByUsername(username string) (*User, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var user User
	err = conn.Users().Find(bson.M{"username": username}).One(&user)
	if err == mgo.ErrNotFound {
		return nil, &errors.ValidationError{Message: "User not found"}
	}

	return &user, nil
}

func (user *User) GetTeams() ([]*Group, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var groups []*Group
	err = conn.Groups().Find(bson.M{"users": bson.M{"$in": []string{user.Username}}}).All(&groups)
	return groups, nil
}
