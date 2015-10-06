package account

import (
	"github.com/apihub/apihub/errors"
	. "github.com/apihub/apihub/log"
	utils "github.com/mrvdot/golang-utils"
)

const (
	ALL_TEAMS string = "*"
)

// The Team type is an encapsulation of a team details.
// It is not allowed to have more than one team with the same alias.
// The `Owner` field indicates the user who created the team.
type Team struct {
	Name     string    `json:"name"`
	Alias    string    `json:"alias"`
	Users    []string  `json:"users"`
	Owner    string    `json:"owner"`
	Services []Service `json:"services,omitempty"`
	Apps     []App     `json:"apps,omitempty"`
}

// Create a team.
//
// It requires to inform the owner and a name.
// If the `alias` is not informed, it will be generate based on the team name.
func (team *Team) Create(owner User) error {
	if err := team.valid(); err != nil {
		Logger.Info("Failed to create a team with invalid data: %+v.", team)
		return err
	}

	team.Users = append(team.Users, owner.Email)
	team.Owner = owner.Email
	if team.Alias == "" {
		team.Alias = utils.GenerateSlug(team.Name)
	} else {
		team.Alias = utils.GenerateSlug(team.Alias)
	}

	if team.Exists() {
		Logger.Info("Failed to create a team with duplicate data: %+v.", team)
		return errors.NewValidationError(errors.ErrTeamDuplicateEntry)
	}

	err := store.UpsertTeam(*team)
	Logger.Info("team.Create: %+v. Err: %s.", team, err)
	return err
}

func (team *Team) Update() error {
	if err := team.valid(); err != nil {
		Logger.Info("Failed to udpate a team with invalid data: %+v.", team)
		return err
	}

	err := store.UpsertTeam(*team)
	Logger.Info("team.Update: %+v. Err: %s.", team, err)
	return err
}

// Delete removes an existing team from the server.
func (team Team) Delete(owner User) error {
	if team.Owner != owner.Email {
		return errors.NewForbiddenError(errors.ErrOnlyOwnerHasPermission)
	}

	go DeleteServicesByTeam(team, owner)
	go DeleteAppsByTeam(team, owner)
	go store.DeleteHooksByTeam(team)

	err := store.DeleteTeam(team)
	Logger.Info("team.Delete: %+v. Err: %s.", team, err)
	return err
}

// Exists checks if there is a team with the same alias in the database.
// Returns `true` if so, and `false` otherwise.
func (team Team) Exists() bool {
	_, err := FindTeamByAlias(team.Alias)
	if err != nil {
		return false
	}
	return true
}

func FindTeamByAlias(alias string) (*Team, error) {
	team, err := store.FindTeamByAlias(alias)
	if err != nil {
		return nil, err
	}

	team.Services, err = store.TeamServices(team)
	if err != nil {
		return nil, err
	}
	team.Apps, err = store.TeamApps(team)
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
	return -1, errors.NewForbiddenError(errors.ErrUserNotInTeam)
}

// Add valid user in the team.
//
// Update the database only if the user is valid.
// Otherwise, ignore invalid or non-existing users.
// Do nothing if the user is already in the team.
func (team *Team) AddUsers(emails []string) error {
	var newUser bool
	var user *User
	for _, email := range emails {
		user = &User{Email: email}
		if !user.Exists() {
			Logger.Info("Failed to add the user '%s' in the team '%s' (User not found).", user.Email, team.Alias)
			continue
		}
		if _, err := team.ContainsUser(user); err != nil {
			team.Users = append(team.Users, user.Email)
			newUser = true
		}
	}

	if newUser {
		err := store.UpsertTeam(*team)
		Logger.Info("team.AddUsers: %+v. Err: %s.", team, err)
		return err
	}
	return nil
}

// Remove a user from the team.
//
// Do nothing if the user is not in the team.
// Return an error if trying to remove the owner. It's not allowed to do that.
func (team *Team) RemoveUsers(emails []string) error {
	var (
		errOwner     errors.ValidationError
		removedUsers bool
		user         *User
		err          interface{}
	)

	for _, email := range emails {
		if team.Owner == email {
			errOwner = errors.NewValidationError(errors.ErrRemoveOwnerFromTeam)
			err = &errOwner
			Logger.Warn("Could not remove the from %s from the team: %s.", team.Owner, team.Alias)
			continue
		}

		user = &User{Email: email}
		if !user.Exists() {
			Logger.Info("Failed to remove the user '%s' from team '%s' (User not found).", user.Email, team.Alias)
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
		err := store.UpsertTeam(*team)
		Logger.Info("team.RemoveUsers: %+v. Err: %s.", team, err)
		return err
	}
	if err != nil {
		return errOwner
	}
	return nil
}

func (team *Team) valid() error {
	if team.Name == "" {
		return errors.NewValidationError(errors.ErrTeamMissingRequiredFields)
	}
	return nil
}
