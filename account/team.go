package account

import (
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
	utils "github.com/mrvdot/golang-utils"
)

// The Team type is an encapsulation of a team details.
// It is not allowed to have more than one team with the same alias.
// The `Owner` field indicates the user who created the team.
type Team struct {
	Name  string   `json:"name"`
	Alias string   `json:"alias"`
	Users []string `json:"users"`
	Owner string   `json:"owner"`
}

// Create a team.
//
// It requires to inform the owner and a name.
// If the `alias` is not informed, it will be generate based on the team name.
func (team *Team) Create(owner User) error {
	if !team.valid() {
		return errors.NewValidationErrorNEW(errors.ErrTeamMissingRequiredFields)
	}

	team.Users = append(team.Users, owner.Email)
	team.Owner = owner.Email
	if team.Alias == "" {
		team.Alias = utils.GenerateSlug(team.Name)
	} else {
		team.Alias = utils.GenerateSlug(team.Alias)
	}

	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	if team.Exists() {
		return errors.NewValidationErrorNEW(errors.ErrTeamDuplicateEntry)
	}

	return store.UpsertTeam(*team)
}

func (team *Team) Update() error {
	if !team.valid() {
		return errors.NewValidationErrorNEW(errors.ErrTeamMissingRequiredFields)
	}

	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	return store.UpsertTeam(*team)
}

// Delete removes an existing team from the server.
func (team Team) Delete(owner User) error {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	if err != nil || team.Owner != owner.Email {
		return errors.NewForbiddenErrorNEW(errors.ErrOnlyOwnerHasPermission)
	}

	err = store.DeleteTeam(team)

	return err
}

// Exists checks if there is a team with the same alias in the database.
// Returns `true` if so, and `false` otherwise.
func (team Team) Exists() bool {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return false
	}
	defer store.Close()

	_, err = store.FindTeamByAlias(team.Alias)
	if err != nil {
		return false
	}
	return true
}

func FindTeamByAlias(alias string) (*Team, error) {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return nil, err
	}
	defer store.Close()

	team, err := store.FindTeamByAlias(alias)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// Check if the user belongs to the team.
// Return the position if so.
func (team *Team) ContainsUser(user *User) (int, error) {
	for i, u := range team.Users {
		if u == user.Email {
			return i, nil
		}
	}
	return -1, errors.NewForbiddenErrorNEW(errors.ErrUserNotInTeam)
}

// Add valid user in the team.
//
// Update the database only if the user is valid.
// Otherwise, ignore invalid or non-existing users.
// Do nothing if the user is already in the team.
func (team *Team) AddUsers(emails []string) error {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	var newUser bool
	var user *User
	for _, email := range emails {
		user = &User{Email: email}
		if !user.Exists() {
			continue
		}
		if _, err := team.ContainsUser(user); err != nil {
			team.Users = append(team.Users, user.Email)
			newUser = true
		}
	}

	if newUser {
		return store.UpsertTeam(*team)
	}
	return nil
}

// Remove a user from the team.
//
// Do nothing if the user is not in the team.
// Return an error if trying to remove the owner. It's not allowed to do that.
func (team *Team) RemoveUsers(emails []string) error {
	store, err := NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer store.Close()

	var (
		errOwner     errors.ValidationErrorNEW
		removedUsers bool
		user         *User
	)
	for _, email := range emails {
		if team.Owner == email {
			errOwner = errors.NewValidationErrorNEW(errors.ErrRemoveOwnerFromTeam)
			err = &errOwner
			continue
		}

		user = &User{Email: email}
		if !user.Exists() {
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
		return store.UpsertTeam(*team)
	}
	if err != nil {
		return errOwner
	}
	return nil
}

func (team *Team) valid() bool {
	if team.Name == "" {
		return false
	}
	return true
}
