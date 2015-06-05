// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
)

type Mem struct {
	Users map[string]account.User
}

func New() account.Storable {
	return &Mem{
		Users: make(map[string]account.User),
	}
}

func (m *Mem) CreateUser(u account.User) error {
	for _, user := range m.Users {
		if user.Email == u.Email || user.Username == u.Username {
			return errors.ErrUserDuplicateEntry
		}
	}

	m.Users[u.Email] = u
	return nil
}

func (m *Mem) UpdateUser(u account.User) error {
	if _, ok := m.Users[u.Email]; !ok {
		return errors.ErrUserNotFound
	}

	m.Users[u.Email] = u
	return nil
}

func (m *Mem) DeleteUser(u account.User) error {
	if _, ok := m.Users[u.Email]; !ok {
		return errors.ErrUserNotFound
	}

	delete(m.Users, u.Email)
	return nil
}

func (m *Mem) FindUserByEmail(email string) (account.User, error) {
	if user, ok := m.Users[email]; !ok {
		return account.User{}, errors.ErrUserNotFound
	} else {
		return user, nil
	}
}

func (m *Mem) Close() {}
