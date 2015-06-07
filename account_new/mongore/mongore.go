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

func (m *Mongore) UpsertUser(u account_new.User) error {
	_, err := m.Users().Upsert(bson.M{"email": u.Email}, u)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteUser(u account_new.User) error {
	err := m.Users().Remove(u)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindUserByEmail(email string) (account_new.User, error) {
	var user account_new.User
	err := m.Users().Find(bson.M{"email": email}).One(&user)

	if err == mgo.ErrNotFound {
		return account_new.User{}, errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return user, err
}

func (m *Mongore) UpsertTeam(t account_new.Team) error {
	_, err := m.Teams().Upsert(bson.M{"alias": t.Alias}, t)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteTeam(t account_new.Team) error {
	err := m.Teams().Remove(t)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindTeamByAlias(alias string) (account_new.Team, error) {
	var team account_new.Team
	err := m.Teams().Find(bson.M{"alias": alias}).One(&team)

	if err == mgo.ErrNotFound {
		return account_new.Team{}, errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return team, err
}

func (m *Mongore) Close() {
	m.store.Close()
}
