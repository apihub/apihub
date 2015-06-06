package account_new

// Storable is an interface for "storage".
// To be compatible, the Storage which implements this interface must pass the acceptance suite that could be found
// in the folder account/test/suite.go.
type Storable interface {
	CreateUser(User) error
	UpdateUser(User) error
	DeleteUser(User) error
	FindUserByEmail(string) (User, error)
	Close()
}

var NewStorable func() (Storable, error)
