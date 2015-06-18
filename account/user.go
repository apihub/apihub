package account

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/backstage/maestro/errors"
	. "github.com/backstage/maestro/log"
)

// The User type is an encapsulation of a user details.
// A valid user is capable to interact with the API to manage teams and services.
type User struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// Create creates a new user account.
//
// It requires to inform the fields: Name, Email and Password.
// It is not allowed to create two users with the same email address.
// It returns an error if the user creation fails.
func (user *User) Create() error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.NewValidationError(errors.ErrUserMissingRequiredFields)
	}

	user.hashPassword()
	if user.Exists() {
		return errors.NewValidationError(errors.ErrUserDuplicateEntry)
	}

	return store.UpsertUser(*user)
}

// Updates the password for an existing account.
func (user *User) ChangePassword() error {
	if !user.Exists() {
		return errors.NewNotFoundError(errors.ErrUserNotFound)
	}

	user.hashPassword()
	return store.UpsertUser(*user)
}

// Delete removes an existing user from the server.
//
// All the teams and services which the corresponding user
// is the only member are deleted along with the user account.
// It returns an error if the user is not found.
func (user User) Delete() error {
	return store.DeleteUser(user)
}

// Exists checks if there is a user with the same email in the database.
// Returns `true` if so, and `false` otherwise.
func (user User) Exists() bool {
	_, err := store.FindUserByEmail(user.Email)
	if err != nil {
		return false
	}
	return true
}

func (user *User) Teams() ([]Team, error) {
	return store.UserTeams(*user)
}

func (user *User) Services() ([]Service, error) {
	return store.UserServices(*user)
}

// Encrypts the user password before saving it in the database.
func (user *User) hashPassword() {
	if hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err == nil {
		user.Password = string(hash[:])
	} else {
		Logger.Error(err.Error())
	}
}
