package app

import (
	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	Name  string
	Users []string
}

func CreateGroup(name string, users []User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	group := Group{
		Name:  name,
		Users: make([]string, len(users)),
	}
	group.Users = getUsernames(users)
	err = conn.Groups().Insert(group)
	if mgo.IsDup(err) {
		message := "Someone already has that group name. Could you try another?"
		return &errors.ValidationError{Message: message}
	}

	return nil
}

func (group *Group) AddUsers(users []User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	var newUser bool
	for _, user := range users {
		if group.containsUser(&user) == false {
			group.Users = append(group.Users, user.Username)
			newUser = true
		}
	}
	if newUser {
		conn.Groups().Update(bson.M{"name": group.Name}, group)
	}
	return nil
}

func DeleteGroupByName(name string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Groups().Remove(bson.M{"name": name})
	if err == mgo.ErrNotFound {
		message := "Group not found."
		return &errors.ValidationError{Message: message}
	}

	return nil
}

func FindGroupByName(name string) (*Group, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var group Group
	err = conn.Groups().Find(bson.M{"name": name}).One(&group)
	if err == mgo.ErrNotFound {
		message := "Group not found."
		return nil, &errors.ValidationError{Message: message}
	}

	return &group, nil
}

func getUsernames(users []User) []string {
	usernames := make([]string, len(users))
	for i, u := range users {
		usernames[i] = u.Username
	}
	return usernames
}

func (group *Group) containsUser(user *User) bool {
	for _, u := range group.Users {
		if u == user.Username {
			return true
		}
	}
	return false
}
