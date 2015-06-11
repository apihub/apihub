// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"fmt"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
)

type Mem struct {
	Users      map[string]account_new.User
	Teams      map[string]account_new.Team
	Tokens     map[string]account_new.TokenInfo
	UserTokens map[string]account_new.User
}

func New() account_new.Storable {
	return &Mem{
		Users:      make(map[string]account_new.User),
		Teams:      make(map[string]account_new.Team),
		Tokens:     make(map[string]account_new.TokenInfo),
		UserTokens: make(map[string]account_new.User),
	}
}

func (m *Mem) UpsertUser(u account_new.User) error {
	m.Users[u.Email] = u
	return nil
}

func (m *Mem) DeleteUser(u account_new.User) error {
	if _, ok := m.Users[u.Email]; !ok {
		return errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}

	delete(m.Users, u.Email)
	return nil
}

func (m *Mem) FindUserByEmail(email string) (account_new.User, error) {
	if user, ok := m.Users[email]; !ok {
		return account_new.User{}, errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	} else {
		return user, nil
	}
}

func (m *Mem) UserTeams(email string) ([]account_new.Team, error) {
	teams := []account_new.Team{}
	for _, team := range m.Teams {
		for _, user := range team.Users {
			if email == user {
				teams = append(teams, team)
			}
		}
	}
	return teams, nil
}

func (m *Mem) UpsertTeam(t account_new.Team) error {
	m.Teams[t.Alias] = t
	return nil
}

func (m *Mem) DeleteTeam(t account_new.Team) error {
	if _, ok := m.Teams[t.Alias]; !ok {
		return errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
	}

	delete(m.Teams, t.Alias)
	return nil
}

func (m *Mem) FindTeamByAlias(alias string) (account_new.Team, error) {
	if team, ok := m.Teams[alias]; !ok {
		return account_new.Team{}, errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
	} else {
		return team, nil
	}
}

func (m *Mem) DeleteTeamByAlias(alias string) error {
	team := account_new.Team{Alias: alias}
	return m.DeleteTeam(team)
}

func (m *Mem) CreateToken(token account_new.TokenInfo) error {
	key := fmt.Sprintf("%s: %s", token.Type, token.User.Email)
	m.Tokens[key] = token
	m.UserTokens[token.Token] = *token.User
	return nil
}

func (m *Mem) DecodeToken(key string, t interface{}) error {
	if token, ok := m.Tokens[key]; ok {
		*t.(*account_new.TokenInfo) = token
	}

	if token, ok := m.UserTokens[key]; ok {
		*t.(*account_new.User) = token
	}
	return nil
}

func (m *Mem) DeleteToken(key string) error {
	delete(m.Tokens, key)
	delete(m.UserTokens, key)
	return nil
}

func (m *Mem) Close() {}
