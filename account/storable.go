package account

// Storable is an interface for "storage".
// To be compatible, the Storage which implements this interface must pass the acceptance suite that could be found
// in the folder account/test/suite.go.
type Storable interface {
	UpsertUser(User) error
	DeleteUser(User) error
	FindUserByEmail(string) (User, error)
	UserTeams(User) ([]Team, error)
	UserServices(User) ([]Service, error)

	UpsertTeam(Team) error
	DeleteTeam(Team) error
	FindTeamByAlias(string) (Team, error)
	DeleteTeamByAlias(string) error

	CreateToken(Token) error
	DeleteToken(key string) error
	DecodeToken(key string, t interface{}) error

	UpsertService(Service) error
	DeleteService(Service) error
	FindServiceBySubdomain(string) (Service, error)
	TeamServices(Team) ([]Service, error)

	UpsertApp(App) error
	FindAppByClientId(string) (App, error)
	DeleteApp(App) error
	TeamApps(Team) ([]App, error)

	UpsertPlugin(Plugin) error
	DeletePlugin(Plugin) error
	DeletePluginsByService(Service) error
	FindPluginByNameAndService(string, Service) (Plugin, error)

	UpsertHook(Hook) error
	DeleteHook(Hook) error
	DeleteHooksByTeam(Team) error
	FindHookByName(string) (Hook, error)
	FindHooksByEventAndTeam(string, string) ([]Hook, error)
	FindHooksByEvent(string) ([]Hook, error)
}

var store Storable

func Storage(s Storable) {
	store = s
}
