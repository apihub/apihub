package account_new

import (
	"encoding/json"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
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
	if user.Name == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		return errors.NewValidationErrorNEW(errors.ErrUserMissingRequiredFields)
	}

	user.hashPassword()
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	err = store.CreateUser(*user)
	return err
}

// Updates the password for an existing account.
func (user *User) ChangePassword() error {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	user.hashPassword()
	err = store.UpdateUser(*user)
	return err
}

// Delete removes an existing user from the server.
//
// All the teams and services which the corresponding user
// is the only member are deleted along with the user account.
// It returns an error if the user is not found.
func (user User) Delete() error {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	err = store.DeleteUser(user)
	return err
}

// Exists checks if the user exists in the database.
// Returns `true` if so, and `false` otherwise.
func (user User) Exists() bool {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return false
	}
	defer store.Close()

	_, err = store.FindUserByEmail(user.Email)
	if err != nil {
		return false
	}
	return true
}

//Return a representation of user but without sensitive data.
func (user User) ToString() string {
	user.Password = ""
	u, _ := json.Marshal(user)
	return string(u)
}

// Encrypts the user password before saving it in the database.
func (user *User) hashPassword() {
	if hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		user.Password = string(hash[:])
	}
}
