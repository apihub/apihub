// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"fmt"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
)

type Mem struct {
	Services   map[string]account.Service
	Users      map[string]account.User
	Teams      map[string]account.Team
	Tokens     map[string]account.TokenInfo
	UserTokens map[string]account.User
}

func New() account.Storable {
	return &Mem{
		Services:   make(map[string]account.Service),
		Users:      make(map[string]account.User),
		Teams:      make(map[string]account.Team),
		Tokens:     make(map[string]account.TokenInfo),
		UserTokens: make(map[string]account.User),
	}
}

func (m *Mem) UpsertUser(u account.User) error {
	m.Users[u.Email] = u
	return nil
}

func (m *Mem) DeleteUser(u account.User) error {
	if _, ok := m.Users[u.Email]; !ok {
		return errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
	}

	delete(m.Users, u.Email)
	return nil
}

func (m *Mem) FindUserByEmail(email string) (account.User, error) {
	if user, ok := m.Users[email]; !ok {
		return account.User{}, errors.NewNotFoundErrorNEW(errors.ErrUserNotFound)
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
		return errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
	}

	delete(m.Teams, t.Alias)
	return nil
}

func (m *Mem) FindTeamByAlias(alias string) (account.Team, error) {
	if team, ok := m.Teams[alias]; !ok {
		return account.Team{}, errors.NewNotFoundErrorNEW(errors.ErrTeamNotFound)
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

func (m *Mem) CreateToken(token account.TokenInfo) error {
	key := fmt.Sprintf("%s: %s", token.Type, token.User.Email)
	m.Tokens[key] = token
	m.UserTokens[token.Token] = *token.User
	return nil
}

func (m *Mem) DecodeToken(key string, t interface{}) error {
	if token, ok := m.Tokens[key]; ok {
		*t.(*account.TokenInfo) = token
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
		return errors.NewNotFoundErrorNEW(errors.ErrServiceNotFound)
	}

	delete(m.Services, s.Subdomain)
	return nil
}

func (m *Mem) FindServiceBySubdomain(subdomain string) (account.Service, error) {
	if service, ok := m.Services[subdomain]; !ok {
		return account.Service{}, errors.NewNotFoundErrorNEW(errors.ErrServiceNotFound)
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

func (m *Mem) Close() {}
