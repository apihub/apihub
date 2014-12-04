package account

import (
	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty""`
	Name  string        `json:"name"`
	Users []string      `json:"users"`
	Owner string        `json:"owner"`
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

func (group *Group) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Groups().Remove(group)
	if err == mgo.ErrNotFound {
		message := "Group not found."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func (group *Group) AddUsers(usernames []string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	var newUser bool
	var user *User
	for _, username := range usernames {
		user = &User{Username: username}
		if !user.Valid() {
			continue
		}
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

func (group *Group) RemoveUsers(usernames []string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	var (
		removedUsers       bool
		errOwnerNotRemoved *errors.ValidationError
		user               *User
	)
	for _, username := range usernames {
		if group.Owner == username {
			message := "It is not possible to remove the owner from the team."
			errOwnerNotRemoved = &errors.ValidationError{Message: message}
			continue
		}

		user = &User{Username: username}
		if !user.Valid() {
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

func FindGroupById(id string) (*Group, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var errNotFound = &errors.ValidationError{Message: "Group not found."}
	if !bson.IsObjectIdHex(id) {
		return nil, errNotFound
	}

	var group Group
	err = conn.Groups().FindId(bson.ObjectIdHex(id)).One(&group)
	if err != nil {
		return nil, errNotFound
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
