package account_new

// Storable is an interface for "storage".
// To be compatible, the Storage which implements this interface must pass the acceptance suite that could be found
// in the folder account/test/suite.go.
type Storable interface {
	// NewStorable() func() (Storable, error)

	UpsertUser(User) error
	DeleteUser(User) error
	FindUserByEmail(string) (User, error)

	UpsertTeam(Team) error
	DeleteTeam(Team) error
	FindTeamByAlias(string) (Team, error)
	Close()
}

var NewStorable func() (Storable, error)
