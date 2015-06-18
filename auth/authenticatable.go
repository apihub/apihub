package auth

import (
	"github.com/backstage/maestro/account"
)

type Authenticatable interface {
	Authenticate(email, password string) (*account.User, bool)
	CreateUserToken(user *account.User) (*account.Token, error)
	UserFromToken(token string) (*account.User, error)
	RevokeUserToken(token string) error
}
