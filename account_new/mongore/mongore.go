package mongore

import (
	"fmt"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongore struct {
	store *storage.Storage
}

func New(config Config) (account_new.Storable, error) {
	conn, err := getConnection(config)
	if err != nil {
		panic(fmt.Sprintf("Error while establishing connection to MongoDB: %s", err.Error()))
		return nil, err
	}

	return &Mongore{
		store: conn,
	}, nil
}

func (m *Mongore) CreateUser(u account_new.User) error {
	err := m.Users().Insert(u)
	if mgo.IsDup(err) {
		Logger.Warn(err.Error())
		return errors.NewValidationErrorNEW(errors.ErrUserDuplicateEntry)
	}

	return err
}

func (m *Mongore) UpdateUser(u account_new.User) error {
	err := m.Users().Update(bson.M{"email": u.Email}, bson.M{"$set": u})
	if err == mgo.ErrNotFound {
		return errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}

	return err
}

func (m *Mongore) DeleteUser(u account_new.User) error {
	err := m.Users().Remove(u)
	if err == mgo.ErrNotFound {
		return errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}

	return err
}

func (m *Mongore) FindUserByEmail(email string) (account_new.User, error) {
	var user account_new.User
	err := m.Users().Find(bson.M{"email": email}).One(&user)
	if err == mgo.ErrNotFound {
		return account_new.User{}, errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}

	return user, err
}

func (m *Mongore) Close() {
	m.store.Close()
}
