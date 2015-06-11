package mongore

import (
	"fmt"
	"time"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
	"github.com/fatih/structs"
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

func (m *Mongore) UserTeams(email string) ([]account_new.Team, error) {
	teams := []account_new.Team{}
	err := m.Teams().Find(bson.M{"users": bson.M{"$in": []string{email}}}).All(&teams)
	return teams, err
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

func (m *Mongore) DeleteTeamByAlias(alias string) error {
	err := m.Teams().Remove(bson.M{"alias": alias})

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) CreateToken(token account_new.TokenInfo) error {
	key := fmt.Sprintf("%s: %s", token.Type, token.User.Email)
	db.Cache.Set(key, nil, time.Duration(token.Expires)*time.Minute)
	db.HMSET(key, token.Expires, structs.Map(token))

	db.Cache.Set(token.Token, nil, time.Duration(token.Expires))
	db.HMSET(token.Token, token.Expires, structs.Map(token.User))
	return nil
}

func (m *Mongore) DecodeToken(key string, t interface{}) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.GetTokenValue(key, t)
}

func (m *Mongore) DeleteToken(key string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.DeleteToken(key)
	return err
}

func (m *Mongore) Close() {
	m.store.Close()
}
