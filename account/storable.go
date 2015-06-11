package account

// Storable is an interface for "storage".
// To be compatible, the Storage which implements this interface must pass the acceptance suite that could be found
// in the folder account/test/suite.go.
type Storable interface {
	UpsertUser(User) error
	DeleteUser(User) error
	FindUserByEmail(string) (User, error)
	UserTeams(string) ([]Team, error)

	UpsertTeam(Team) error
	DeleteTeam(Team) error
	FindTeamByAlias(string) (Team, error)
	DeleteTeamByAlias(string) error

	CreateToken(TokenInfo) error
	DeleteToken(string) error
	DecodeToken(key string, t interface{}) error
	Close()
}

var NewStorable func() (Storable, error)
