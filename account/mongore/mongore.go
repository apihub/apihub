package mongore

import (
	"fmt"
	"time"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/db"
	"github.com/backstage/maestro/errors"
	. "github.com/backstage/maestro/log"
	"github.com/fatih/structs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongore struct {
	config Config
}

func New(config Config) account.Storable {
	return &Mongore{config: config}
}

func (m *Mongore) UpsertUser(u account.User) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	_, err := strg.Users().Upsert(bson.M{"email": u.Email}, u)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteUser(u account.User) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Users().Remove(u)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrUserNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindUserByEmail(email string) (account.User, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	user := account.User{}
	err := strg.Users().Find(bson.M{"email": email}).One(&user)

	if err == mgo.ErrNotFound {
		return account.User{}, errors.NewNotFoundError(errors.ErrUserNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return user, err
}

func (m *Mongore) UserTeams(user account.User) ([]account.Team, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	teams := []account.Team{}
	err := strg.Teams().Find(bson.M{"users": bson.M{"$in": []string{user.Email}}}).All(&teams)
	return teams, err
}

func (m *Mongore) UpsertTeam(t account.Team) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	_, err := strg.Teams().Upsert(bson.M{"alias": t.Alias}, t)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteTeam(t account.Team) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Teams().Remove(t)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrTeamNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindTeamByAlias(alias string) (account.Team, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	team := account.Team{}
	err := strg.Teams().Find(bson.M{"alias": alias}).One(&team)

	if err == mgo.ErrNotFound {
		return account.Team{}, errors.NewNotFoundError(errors.ErrTeamNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return team, err
}

func (m *Mongore) DeleteTeamByAlias(alias string) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Teams().Remove(bson.M{"alias": alias})

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrTeamNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) TeamServices(team account.Team) ([]account.Service, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	services := []account.Service{}
	err := strg.Services().Find(bson.M{"team": team.Alias}).All(&services)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return services, err
}

func (m *Mongore) CreateToken(token account.Token) error {
	key := fmt.Sprintf("%s: %s", token.Type, token.User.Email)
	db.Cache.Set(key, nil, time.Duration(token.Expires)*time.Minute)
	db.HMSET(key, token.Expires, structs.Map(token))

	db.Cache.Set(token.AccessToken, nil, time.Duration(token.Expires))
	db.HMSET(token.AccessToken, token.Expires, structs.Map(token.User))
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

func (m *Mongore) UpsertService(s account.Service) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	_, err := strg.Services().Upsert(bson.M{"subdomain": s.Subdomain}, s)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteService(s account.Service) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Services().Remove(s)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrServiceNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindServiceBySubdomain(subdomain string) (account.Service, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	var service account.Service
	err := strg.Services().Find(bson.M{"subdomain": subdomain}).One(&service)

	if err == mgo.ErrNotFound {
		return account.Service{}, errors.NewNotFoundError(errors.ErrServiceNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return service, err
}

func (m *Mongore) UserServices(user account.User) ([]account.Service, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	var services []account.Service = []account.Service{}

	teams, err := m.UserTeams(user)
	if err != nil {
		Logger.Warn(err.Error())
	}
	var st []string = make([]string, len(teams))
	for i, team := range teams {
		st[i] = team.Alias
	}

	err = strg.Services().Find(bson.M{"team": bson.M{"$in": st}}).All(&services)
	return services, err
}

func (m *Mongore) UpsertApp(app account.App) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	_, err := strg.Apps().Upsert(bson.M{"clientid": app.ClientId}, app)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindAppByClientId(clientid string) (account.App, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	var app account.App
	err := strg.Apps().Find(bson.M{"clientid": clientid}).One(&app)

	if err == mgo.ErrNotFound {
		return account.App{}, errors.NewNotFoundError(errors.ErrAppNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return app, err
}

func (m *Mongore) DeleteApp(app account.App) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Apps().Remove(bson.M{"clientid": app.ClientId})

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrAppNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) TeamApps(team account.Team) ([]account.App, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	apps := []account.App{}
	err := strg.Apps().Find(bson.M{"team": team.Alias}).All(&apps)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return apps, err
}

func (m *Mongore) UpsertPlugin(pc account.Plugin) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	_, err := strg.PluginsConfig().Upsert(bson.M{"service": pc.Service, "name": pc.Name}, pc)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeletePlugin(pc account.Plugin) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.PluginsConfig().Remove(pc)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrPluginNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeletePluginsByService(service account.Service) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.PluginsConfig().Remove(bson.M{"service": service.Subdomain})

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrPluginNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindPluginByNameAndService(pluginName string, service account.Service) (account.Plugin, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	var plugin account.Plugin
	err := strg.PluginsConfig().Find(bson.M{"name": pluginName, "service": service.Subdomain}).One(&plugin)

	if err == mgo.ErrNotFound {
		return account.Plugin{}, errors.NewNotFoundError(errors.ErrPluginNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return plugin, err
}

func (m *Mongore) UpsertHook(w account.Hook) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	_, err := strg.Hooks().Upsert(bson.M{"name": w.Name}, w)

	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteHook(w account.Hook) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Hooks().Remove(w)

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrHookNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) DeleteHooksByTeam(team account.Team) error {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	err := strg.Hooks().Remove(bson.M{"team": team.Alias})

	if err == mgo.ErrNotFound {
		return errors.NewNotFoundError(errors.ErrHookNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return err
}

func (m *Mongore) FindHookByName(name string) (account.Hook, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	hooks := account.Hook{}
	err := strg.Hooks().Find(bson.M{"name": name}).One(&hooks)

	if err == mgo.ErrNotFound {
		return account.Hook{}, errors.NewNotFoundError(errors.ErrHookNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return hooks, err
}

func (m *Mongore) FindHooksByEvent(event string) ([]account.Hook, error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	hooks := []account.Hook{}
	err := strg.Hooks().Find(bson.M{"events": bson.M{"$in": []string{event}}}).All(&hooks)

	if err == mgo.ErrNotFound {
		return []account.Hook{}, errors.NewNotFoundError(errors.ErrHookNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return hooks, err
}

func (m *Mongore) FindHooksByEventAndTeam(event string, team string) (hooks []account.Hook, err error) {
	var strg Storage
	strg.Storage = m.openSession()
	defer strg.Close()

	hooks = []account.Hook{}
	if team == account.ALL_TEAMS {
		err = strg.Hooks().Find(bson.M{"events": bson.M{"$in": []string{event}}}).All(&hooks)
	} else {
		err = strg.Hooks().Find(bson.M{"team": team, "events": bson.M{"$in": []string{event}}}).All(&hooks)
	}

	if err == mgo.ErrNotFound {
		return []account.Hook{}, errors.NewNotFoundError(errors.ErrHookNotFound)
	}
	if err != nil {
		Logger.Warn(err.Error())
	}

	return hooks, err
}
