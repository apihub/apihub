package auth

import (
	"github.com/backstage/apimanager/account"
)

type Authenticatable interface {
	Authenticate(email, password string) (*account.User, bool)
	CreateUserToken(*account.User) (*account.Token, error)
	UserFromToken(token string) (*account.User, error)
	RevokeUserToken(token string) error
}
