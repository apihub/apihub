package account

import (
	"github.com/albertoleal/backstage/db"
	"github.com/albertoleal/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Team struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty""`
	Name  string        `json:"name"`
	Users []string      `json:"users"`
	Owner string        `json:"owner"`
}

func (team *Team) Save(owner *User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	team.Users = []string{owner.Username}
	team.Owner = owner.Username
	err = conn.Teams().Insert(team)
	if mgo.IsDup(err) {
		message := "Someone already has that team name. Could you try another?"
		return &errors.ValidationError{Message: message}
	}

	return nil
}

func (team *Team) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Teams().Remove(team)
	if err == mgo.ErrNotFound {
		message := "Team not found."
		return &errors.ValidationError{Message: message}
	}
	return err
}

func (team *Team) AddUsers(usernames []string) error {
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
		if _, contains := team.ContainsUser(user); contains == false {
			team.Users = append(team.Users, user.Username)
			newUser = true
		}
	}
	if newUser {
		conn.Teams().Update(bson.M{"name": team.Name}, team)
	}
	return nil
}

func (team *Team) RemoveUsers(usernames []string) error {
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
		if team.Owner == username {
			message := "It is not possible to remove the owner from the team."
			errOwnerNotRemoved = &errors.ValidationError{Message: message}
			continue
		}

		user = &User{Username: username}
		if !user.Valid() {
			continue
		}
		if i, ok := team.ContainsUser(user); ok {
			hi := len(team.Users) - 1
			if hi > i {
				team.Users[i] = team.Users[hi]
			}
			team.Users = team.Users[:hi]
			removedUsers = true
		}
	}
	if removedUsers {
		conn.Teams().Update(bson.M{"name": team.Name}, team)
	}
	if errOwnerNotRemoved != nil {
		return errOwnerNotRemoved
	}
	return nil
}

func DeleteTeamByName(name string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Teams().Remove(bson.M{"name": name})
	if err == mgo.ErrNotFound {
		message := "Team not found."
		return &errors.ValidationError{Message: message}
	}

	return nil
}

func FindTeamByName(name string) (*Team, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var team Team
	err = conn.Teams().Find(bson.M{"name": name}).One(&team)
	if err == mgo.ErrNotFound {
		message := "Team not found."
		return nil, &errors.ValidationError{Message: message}
	}

	return &team, nil
}

func FindTeamById(id string) (*Team, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var errNotFound = &errors.ValidationError{Message: "Team not found."}
	if !bson.IsObjectIdHex(id) {
		return nil, errNotFound
	}

	var team Team
	err = conn.Teams().FindId(bson.ObjectIdHex(id)).One(&team)
	if err != nil {
		return nil, errNotFound
	}

	return &team, nil
}

func (team *Team) GetTeamUsers() ([]*User, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var users []*User
	var user *User
	for _, username := range team.Users {
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

func (team *Team) ContainsUser(user *User) (int, bool) {
	for i, u := range team.Users {
		if u == user.Username {
			return i, true
		}
	}
	return -1, false
}
