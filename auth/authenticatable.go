package auth

import (
	"github.com/apihub/apihub/account"
)

type Authenticatable interface {
	Authenticate(email, password string) (*account.User, bool)
	// ChangePassword(email, password string) (*account.User, bool)
	// RecoverPassword(email string) (bool, error)
	CreateUserToken(user *account.User) (*account.Token, error)
	UserFromToken(token string) (*account.User, error)
	RevokeUserToken(token string) error
}
