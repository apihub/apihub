// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
)

type Mem struct {
	Users map[string]account_new.User
	Teams map[string]account_new.Team
}

func New() account_new.Storable {
	return &Mem{
		Users: make(map[string]account_new.User),
		Teams: make(map[string]account_new.Team),
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

func (m *Mem) Close() {}
