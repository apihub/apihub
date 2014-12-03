package account

import (
	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	Owner string   `json:"owner"`
}

func (group *Group) Save(owner *User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	group.Users = []string{owner.Username}
	group.Owner = owner.Username
	err = conn.Groups().Insert(group)
	if mgo.IsDup(err) {
		message := "Someone already has that group name. Could you try another?"
		return &errors.ValidationError{Message: message}
	}

	return nil
}

func (group *Group) AddUsers(users []*User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	var newUser bool
	for _, user := range users {
		if _, contains := group.containsUser(user); contains == false {
			group.Users = append(group.Users, user.Username)
			newUser = true
		}
	}
	if newUser {
		conn.Groups().Update(bson.M{"name": group.Name}, group)
	}
	return nil
}

func (group *Group) RemoveUsers(users []*User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	var (
		removedUsers       bool
		errOwnerNotRemoved *errors.ValidationError
	)
	for _, user := range users {
		if group.Owner == user.Username {
			message := "It is not possible to remove the owner from the team."
			errOwnerNotRemoved = &errors.ValidationError{Message: message}
			continue
		}

		if i, ok := group.containsUser(user); ok {
			hi := len(group.Users) - 1
			if hi > i {
				group.Users[i] = group.Users[hi]
			}
			group.Users = group.Users[:hi]
			removedUsers = true
		}
	}
	if removedUsers {
		conn.Groups().Update(bson.M{"name": group.Name}, group)
	}
	if errOwnerNotRemoved != nil {
		return errOwnerNotRemoved
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

func (group *Group) GetGroupUsers() ([]*User, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var users []*User
	var user *User
	for _, username := range group.Users {
		user, _ = FindUserByUsername(username)
		users = append(users, user)
	}

	return users, nil
}

func getUsernames(users []*User) []string {
	usernames := make([]string, len(users))
	for i, u := range users {
		usernames[i] = u.Username
	}
	return usernames
}

func (group *Group) containsUser(user *User) (int, bool) {
	for i, u := range group.Users {
		if u == user.Username {
			return i, true
		}
	}
	return -1, false
}
