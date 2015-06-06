// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
)

type Mem struct {
	Users map[string]account_new.User
}

func New() account_new.Storable {
	return &Mem{
		Users: make(map[string]account_new.User),
	}
}

func (m *Mem) CreateUser(u account_new.User) error {
	for _, user := range m.Users {
		if user.Email == u.Email || user.Username == u.Username {
			return errors.NewValidationErrorNEW(errors.ErrUserDuplicateEntry)
		}
	}

	m.Users[u.Email] = u
	return nil
}

func (m *Mem) UpdateUser(u account_new.User) error {
	if _, ok := m.Users[u.Email]; !ok {
		return errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}

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

func (m *Mem) Close() {}
