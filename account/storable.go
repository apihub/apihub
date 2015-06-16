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

	CreateToken(TokenInfo) error
	DeleteToken(key string) error
	DecodeToken(key string, t interface{}) error

	UpsertService(Service) error
	DeleteService(Service) error
	FindServiceBySubdomain(string) (Service, error)
	TeamServices(Team) ([]Service, error)

	UpsertApp(App) error
	FindAppByClientId(string) (App, error)
	DeleteApp(App) error

	UpsertPluginConfig(PluginConfig) error
	DeletePluginConfig(PluginConfig) error
	FindPluginConfigByNameAndService(string, Service) (PluginConfig, error)
}

var store Storable

func Storage(s Storable) {
	store = s
}
