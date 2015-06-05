package account

import (
	"encoding/json"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// The User type is an encapsulation of a user details.
// A valid user is capable to interact with the API to manage teams and services.
type User struct {
	Name                 string `json:"name,omitempty"`
	Email                string `json:"email,omitempty"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	NewPassword          string `json:"new_password,omitempty" bson:"-"`
	ConfirmationPassword string `json:"confirmation_password,omitempty" bson:"-"`
}

// Save creates a new user account.
//
// It requires to inform the fields: Name, Email and Password.
// It is not allowed to create two users with the same email address.
// It returns an error if the user creation fails.
func (user *User) Save() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if user.Name == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		return &errors.ValidationError{Payload: "Name/Email/Username/Password cannot be empty."}
	}

	user.HashPassword()
	err = conn.Users().Insert(user)
	if mgo.IsDup(err) {
		return &errors.ValidationError{Payload: "Someone already has that email/username. Could you try another?"}
	}
	return err
}

// Updates the password for an existing account
func (user *User) ChangePassword() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	user.HashPassword()
	err = conn.Users().Update(bson.M{"email": user.Email}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return err
}

// Delete removes an existing user from the server.
//
// All the teams and services which the corresponding user
// is the only member are deleted along with the user account.
// It returns an error if the user is not found.
func (user *User) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = user.remove()
	if err == mgo.ErrNotFound {
		return &errors.ValidationError{Payload: "User not found."}
	}
	return err
}

func (user *User) remove() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	var ts []*Team = []*Team{}
	err = conn.Teams().Find(bson.M{"users": bson.M{"$size": 1}, "owner": user.Email}).All(&ts)
	if err != nil {
		return err
	}
	var teams []string
	for _, t := range ts {
		teams = append(teams, t.Alias)
		DeleteServicesByTeam(t.Alias)
		DeleteClientByTeam(t.Alias)
	}
	_, err = conn.Teams().RemoveAll(bson.M{"alias": bson.M{"$in": teams}})
	if err != nil {
		return err
	}

	err = conn.Users().Remove(user)
	if err != nil {
		return err
	}
	return nil
}

// Encrypts the user password before saving it in the database.
func (user *User) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hash[:])
}

// Exists checks if the user exists in the database.
// Returns `true` if so, and `false` otherwise.
func (user *User) Exists() bool {
	_, err := FindUserByEmail(user.Email)
	if err != nil {
		return false
	}
	return true
}

// Try to find a user by its email address.
// If the user is not found, return an error. Return the user otherwise.
func FindUserByEmail(email string) (*User, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var user User
	err = conn.Users().Find(bson.M{"email": email}).One(&user)
	if err == mgo.ErrNotFound {
		return nil, &errors.ValidationError{Payload: "User not found"}
	}
	return &user, nil
}

// Return a list of all the teams which the user belongs to.
func (user *User) GetTeams() ([]*Team, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var teams []*Team = []*Team{}
	err = conn.Teams().Find(bson.M{"users": bson.M{"$in": []string{user.Email}}}).All(&teams)
	return teams, nil
}

// Return a list of all the services which the user belongs to.
func (user *User) GetServices() ([]*Service, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	teams, _ := user.GetTeams()
	var st []string = make([]string, len(teams))
	for i, team := range teams {
		st[i] = team.Alias
	}
	var services []*Service = []*Service{}
	err = conn.Services().Find(bson.M{"team": bson.M{"$in": st}}).All(&services)
	return services, nil
}

//Return a representation of user but without sensitive data.
func (user *User) ToString() string {
	user.Password = ""
	u, _ := json.Marshal(user)
	return string(u)
}
