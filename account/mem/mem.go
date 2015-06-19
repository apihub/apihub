// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"fmt"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
)

type Mem struct {
	Apps          map[string]account.App
	Services      map[string]account.Service
	Users         map[string]account.User
	Teams         map[string]account.Team
	PluginsConfig map[string]map[string]account.PluginConfig
	Tokens        map[string]account.Token
	UserTokens    map[string]account.User
	Webhooks      map[string]account.Webhook
}

func New() account.Storable {
	return &Mem{
		Apps:          make(map[string]account.App),
		Services:      make(map[string]account.Service),
		Users:         make(map[string]account.User),
		Teams:         make(map[string]account.Team),
		PluginsConfig: make(map[string]map[string]account.PluginConfig),
		Tokens:        make(map[string]account.Token),
		UserTokens:    make(map[string]account.User),
		Webhooks:      make(map[string]account.Webhook),
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

func (m *Mem) UpsertPluginConfig(pc account.PluginConfig) error {
	m.PluginsConfig[pc.Service] = map[string]account.PluginConfig{pc.Name: pc}
	return nil
}

func (m *Mem) DeletePluginConfig(pc account.PluginConfig) error {
	if _, ok := m.PluginsConfig[pc.Service][pc.Name]; !ok {
		return errors.NewNotFoundError(errors.ErrPluginConfigNotFound)
	}

	delete(m.PluginsConfig, pc.Name)
	return nil
}
func (m *Mem) FindPluginConfigByNameAndService(pluginName string, service account.Service) (account.PluginConfig, error) {
	if plugin, ok := m.PluginsConfig[service.Subdomain][pluginName]; !ok {
		return account.PluginConfig{}, errors.NewNotFoundError(errors.ErrPluginConfigNotFound)
	} else {
		return plugin, nil
	}
}

func (m *Mem) UpsertWebhook(w account.Webhook) error {
	m.Webhooks[w.Name] = w
	return nil
}

func (m *Mem) DeleteWebhook(w account.Webhook) error {
	if _, ok := m.Webhooks[w.Name]; !ok {
		return errors.NewNotFoundError(errors.ErrWebhookNotFound)
	}

	delete(m.Webhooks, w.Name)
	return nil
}

func (m *Mem) FindWebhookByName(name string) (account.Webhook, error) {
	if webhook, ok := m.Webhooks[name]; !ok {
		return account.Webhook{}, errors.NewNotFoundError(errors.ErrWebhookNotFound)
	} else {
		return webhook, nil
	}
}

func (m *Mem) FindWebhooksByEventAndTeam(event string, team string) ([]account.Webhook, error) {
	whs := []account.Webhook{}

	for _, wh := range m.Webhooks {
		for _, ev := range wh.Events {
			if ev == event && (team == "*" || wh.Team == team) {
				whs = append(whs, wh)
			}
		}
	}

	return whs, nil
}
