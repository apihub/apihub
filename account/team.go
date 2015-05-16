package account

import (
	"encoding/json"

	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	. "github.com/mrvdot/golang-utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// The Team type is an encapsulation of a team details.
// It is not allowed to have more than one team with the same alias.
// The `Owner` field indicates the user who created the team.
type Team struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"-""`
	Name     string        `json:"name"`
	Alias    string        `json:"alias"`
	Users    []string      `json:"users"`
	Owner    string        `json:"owner"`
	Services []*Service    `bson:"-" json:"services,omitempty"`
	Clients  []*Client     `bson:"-" json:"clients,omitempty"`
}

// Save creates a new team.
//
// It requires to inform the owner and a name.
// If the `alias` is not informed, it will be generate based on the team name.
func (team *Team) Save(owner *User) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if team.Name == "" {
		return &errors.ValidationError{Payload: "Name cannot be empty."}
	}

	team.Users = []string{owner.Email}
	team.Owner = owner.Email
	if team.Alias == "" {
		team.Alias = GenerateSlug(team.Name)
	} else {
		team.Alias = GenerateSlug(team.Alias)
	}
	err = conn.Teams().Insert(team)
	if mgo.IsDup(err) {
		return &errors.ValidationError{Payload: "Someone already has that team alias. Could you try another?"}
	}

	return nil
}

// DeleteTeamByAlias removes an existing team from the server based on given alias.
//
// Only the owner is allowed to remove the team. If the user is not the owner,
// an error will be returned.
// Deletes all the services that belong to the team.
func DeleteTeamByAlias(alias string, user *User) (*Team, error) {
	team, err := FindTeamByAlias(alias, user)
	if err != nil || team.Owner != user.Email {
		return nil, &errors.ForbiddenError{Payload: errors.ErrOnlyOwnerHasPermission.Error()}
	}
	team.Services, err = FindServicesByTeam(alias)
	if err != nil {
		return nil, err
	}
	team.Clients, err = FindClientsByTeam(alias)
	if err != nil {
		return nil, err
	}
	return team, team.delete()
}

func (team *Team) delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Teams().RemoveId(team.Id)
	if err == mgo.ErrNotFound {
		return &errors.ValidationError{Payload: "Team not found."}
	}

	DeleteServicesByTeam(team.Alias)
	DeleteClientByTeam(team.Alias)
	return nil
}

// Add valid user in the team.
//
// Update the database only if the user is valid.
// Otherwise, ignore invalid or non-existing users.
// Do nothing if the user is already in the team.
func (team *Team) AddUsers(emails []string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	var newUser bool
	var user *User
	for _, email := range emails {
		user = &User{Email: email}
		if !user.Valid() {
			continue
		}
		if _, err := team.ContainsUser(user); err != nil {
			team.Users = append(team.Users, user.Email)
			newUser = true
		}
	}
	if newUser {
		conn.Teams().Update(bson.M{"name": team.Name}, team)
	}
	return nil
}

// Remove a user from the team.
//
// Do nothing if the user is not in the team.
// Return an error if trying to remove the owner. It's not allowed to do that.
func (team *Team) RemoveUsers(emails []string) error {
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
	for _, email := range emails {
		if team.Owner == email {
			errOwnerNotRemoved = &errors.ValidationError{Payload: "It is not possible to remove the owner from the team."}
			continue
		}

		user = &User{Email: email}
		if !user.Valid() {
			continue
		}
		if i, err := team.ContainsUser(user); err == nil {
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

// DeleteTeamByName removes an existing team from the server based on given name.
// Unlike the `Delete` method, it does not delete the services. Be aware of this.
func DeleteTeamByName(name string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Teams().Remove(bson.M{"name": name})
	if err == mgo.ErrNotFound {
		return &errors.ValidationError{Payload: "Team not found."}
	}

	return nil
}

// Find the team info and all the services for a given team name.
func FindTeamByName(name string) (*Team, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var team Team
	err = conn.Teams().Find(bson.M{"name": name}).One(&team)
	if err == mgo.ErrNotFound {
		return nil, &errors.ValidationError{Payload: "Team not found."}
	}
	team.Services, err = FindServicesByTeam(team.Alias)
	if err != nil {
		return nil, err
	}
	team.Clients, err = FindClientsByTeam(team.Alias)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// Find the team info, clients and all the services for a given team alias.
// It returns the team info if the user belongs to the team.
// Return an error otherwise.
func FindTeamByAlias(alias string, user *User) (*Team, error) {
	team, err := findTeamByAlias(alias)
	if err != nil {
		return nil, &errors.NotFoundError{Payload: err.Error()}
	}
	_, err = team.ContainsUser(user)
	if err != nil {
		return nil, &errors.ForbiddenError{Payload: err.Error()}
	}
	team.Services, err = FindServicesByTeam(alias)
	if err != nil {
		return nil, err
	}
	team.Clients, err = FindClientsByTeam(alias)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func findTeamByAlias(alias string) (*Team, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var team Team
	err = conn.Teams().Find(bson.M{"alias": alias}).One(&team)
	if err == mgo.ErrNotFound {
		return nil, &errors.ValidationError{Payload: "Team not found."}
	}

	return &team, nil
}

// Find the team info, clients and all the services for a given team id.
// Unlike the `FindTeamByAlias` method, it does not check if the
// user belong to the team.
func FindTeamById(id string) (*Team, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var errNotFound = &errors.ValidationError{Payload: "Team not found."}
	if !bson.IsObjectIdHex(id) {
		return nil, errNotFound
	}

	var team Team
	err = conn.Teams().FindId(bson.ObjectIdHex(id)).One(&team)
	if err != nil {
		return nil, errNotFound
	}

	team.Services, err = FindServicesByTeam(team.Alias)
	if err != nil {
		return nil, err
	}
	team.Clients, err = FindClientsByTeam(team.Alias)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// Return a list of users that belongs to the given team.
func (team *Team) GetTeamUsers() ([]*User, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var users []*User
	var user *User
	for _, email := range team.Users {
		user, _ = FindUserByEmail(email)
		users = append(users, user)
	}

	return users, nil
}

func getEmails(users []*User) []string {
	emails := make([]string, len(users))
	for i, u := range users {
		emails[i] = u.Email
	}
	return emails
}

// Check if the user belongs to the team.
// Return the position if so.
func (team *Team) ContainsUser(user *User) (int, error) {
	for i, u := range team.Users {
		if u == user.Email {
			return i, nil
		}
	}
	return -1, errors.ErrUserNotInTeam
}

//Return a representation of a team without sensitive data.
func (team *Team) ToString() string {
	team.Id = ""
	t, _ := json.Marshal(team)
	return string(t)
}
