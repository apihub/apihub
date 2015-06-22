// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"fmt"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
)

type Mem struct {
	Apps       map[string]account.App
	Services   map[string]account.Service
	Users      map[string]account.User
	Teams      map[string]account.Team
	Plugins    map[string]map[string]account.Plugin
	Tokens     map[string]account.Token
	UserTokens map[string]account.User
	Hooks      map[string]account.Hook
}

func New() account.Storable {
	return &Mem{
		Apps:       make(map[string]account.App),
		Services:   make(map[string]account.Service),
		Users:      make(map[string]account.User),
		Teams:      make(map[string]account.Team),
		Plugins:    make(map[string]map[string]account.Plugin),
		Tokens:     make(map[string]account.Token),
		UserTokens: make(map[string]account.User),
		Hooks:      make(map[string]account.Hook),
	}
}

func (m *Mem) UpsertUser(u account.User) error {
	m.Users[u.Email] = u
	return nil
}

func (m *Mem) DeleteUser(u account.User) error {
	if _, ok := m.Users[u.Email]; !ok {
		return errors.NewNotFoundError(errors.ErrUserNotFound)
	}

	delete(m.Users, u.Email)
	return nil
}

func (m *Mem) FindUserByEmail(email string) (account.User, error) {
	if user, ok := m.Users[email]; !ok {
		return account.User{}, errors.NewNotFoundError(errors.ErrUserNotFound)
	} else {
		return user, nil
	}
}

func (m *Mem) UserTeams(user account.User) ([]account.Team, error) {
	teams := []account.Team{}
	for _, team := range m.Teams {
		for _, u := range team.Users {
			if user.Email == u {
				teams = append(teams, team)
			}
		}
	}
	return teams, nil
}

func (m *Mem) UpsertTeam(t account.Team) error {
	m.Teams[t.Alias] = t
	return nil
}

func (m *Mem) DeleteTeam(t account.Team) error {
	if _, ok := m.Teams[t.Alias]; !ok {
		return errors.NewNotFoundError(errors.ErrTeamNotFound)
	}

	delete(m.Teams, t.Alias)
	return nil
}

func (m *Mem) FindTeamByAlias(alias string) (account.Team, error) {
	if team, ok := m.Teams[alias]; !ok {
		return account.Team{}, errors.NewNotFoundError(errors.ErrTeamNotFound)
	} else {
		return team, nil
	}
}

func (m *Mem) DeleteTeamByAlias(alias string) error {
	team := account.Team{Alias: alias}
	return m.DeleteTeam(team)
}

func (m *Mem) TeamServices(team account.Team) ([]account.Service, error) {
	services := []account.Service{}
	for _, service := range m.Services {
		if service.Team == team.Alias {
			services = append(services, service)
		}
	}
	return services, nil
}

func (m *Mem) TeamApps(team account.Team) ([]account.App, error) {
	apps := []account.App{}
	for _, app := range m.Apps {
		if app.Team == team.Alias {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

func (m *Mem) CreateToken(token account.Token) error {
	key := fmt.Sprintf("%s: %s", token.Type, token.User.Email)
	m.Tokens[key] = token
	m.UserTokens[token.AccessToken] = *token.User
	return nil
}

func (m *Mem) DecodeToken(key string, t interface{}) error {
	if token, ok := m.Tokens[key]; ok {
		*t.(*account.Token) = token
	}

	if token, ok := m.UserTokens[key]; ok {
		*t.(*account.User) = token
	}
	return nil
}

func (m *Mem) DeleteToken(key string) error {
	delete(m.Tokens, key)
	delete(m.UserTokens, key)
	return nil
}

func (m *Mem) UpsertService(s account.Service) error {
	m.Services[s.Subdomain] = s
	return nil
}

func (m *Mem) DeleteService(s account.Service) error {
	if _, ok := m.Services[s.Subdomain]; !ok {
		return errors.NewNotFoundError(errors.ErrServiceNotFound)
	}

	delete(m.Services, s.Subdomain)
	return nil
}

func (m *Mem) FindServiceBySubdomain(subdomain string) (account.Service, error) {
	if service, ok := m.Services[subdomain]; !ok {
		return account.Service{}, errors.NewNotFoundError(errors.ErrServiceNotFound)
	} else {
		return service, nil
	}
}

func (m *Mem) UserServices(user account.User) ([]account.Service, error) {
	teams, _ := m.UserTeams(user)
	services := []account.Service{}

	var teamServices []account.Service
	for _, team := range teams {
		teamServices, _ = m.TeamServices(team)
		if len(teamServices) > 0 {
			services = append(services, teamServices...)
		}
	}
	return services, nil
}

func (m *Mem) UpsertApp(a account.App) error {
	m.Apps[a.ClientId] = a
	return nil
}

func (m *Mem) FindAppByClientId(id string) (account.App, error) {
	if app, ok := m.Apps[id]; !ok {
		return account.App{}, errors.NewNotFoundError(errors.ErrAppNotFound)
	} else {
		return app, nil
	}
}

func (m *Mem) DeleteApp(a account.App) error {
	if _, ok := m.Apps[a.ClientId]; !ok {
		return errors.NewNotFoundError(errors.ErrAppNotFound)
	}

	delete(m.Apps, a.ClientId)
	return nil
}

func (m *Mem) UpsertPlugin(pc account.Plugin) error {
	m.Plugins[pc.Service] = map[string]account.Plugin{pc.Name: pc}
	return nil
}

func (m *Mem) DeletePlugin(pc account.Plugin) error {
	if _, ok := m.Plugins[pc.Service][pc.Name]; !ok {
		return errors.NewNotFoundError(errors.ErrPluginNotFound)
	}

	delete(m.Plugins, pc.Name)
	return nil
}

func (m *Mem) DeletePluginsByService(service account.Service) error {
	if _, ok := m.Plugins[service.Subdomain]; !ok {
		return errors.NewNotFoundError(errors.ErrPluginNotFound)
	}

	delete(m.Plugins, service.Subdomain)
	return nil
}

func (m *Mem) FindPluginByNameAndService(pluginName string, service account.Service) (account.Plugin, error) {
	if plugin, ok := m.Plugins[service.Subdomain][pluginName]; !ok {
		return account.Plugin{}, errors.NewNotFoundError(errors.ErrPluginNotFound)
	} else {
		return plugin, nil
	}
}

func (m *Mem) UpsertHook(w account.Hook) error {
	m.Hooks[w.Name] = w
	return nil
}

func (m *Mem) DeleteHook(w account.Hook) error {
	if _, ok := m.Hooks[w.Name]; !ok {
		return errors.NewNotFoundError(errors.ErrHookNotFound)
	}

	delete(m.Hooks, w.Name)
	return nil
}

func (m *Mem) DeleteHooksByTeam(team account.Team) error {
	found := false
	for _, wh := range m.Hooks {
		if wh.Team == team.Alias {
			delete(m.Hooks, wh.Name)
			found = true
		}
	}
	if !found {
		return errors.NewNotFoundError(errors.ErrTeamNotFound)
	}
	return nil
}

func (m *Mem) FindHookByName(name string) (account.Hook, error) {
	if hook, ok := m.Hooks[name]; !ok {
		return account.Hook{}, errors.NewNotFoundError(errors.ErrHookNotFound)
	} else {
		return hook, nil
	}
}

func (m *Mem) FindHooksByEvent(event string) ([]account.Hook, error) {
	whs := []account.Hook{}

	for _, wh := range m.Hooks {
		for _, ev := range wh.Events {
			if ev == event {
				whs = append(whs, wh)
			}
		}
	}

	return whs, nil
}

func (m *Mem) FindHooksByEventAndTeam(event string, team string) ([]account.Hook, error) {
	whs := []account.Hook{}

	for _, wh := range m.Hooks {
		for _, ev := range wh.Events {
			if ev == event && (team == account.ALL_TEAMS || wh.Team == team) {
				whs = append(whs, wh)
			}
		}
	}

	return whs, nil
}
