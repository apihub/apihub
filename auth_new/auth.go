package auth_new

import (
	"fmt"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
)

type Authenticatable interface {
	Authenticate(email, password string) (*account_new.User, bool)
	CreateUserToken(*account_new.User) (*ApiToken, error)
	UserFromToken(token string) (*account_new.User, error)
	RevokeUserToken(token string)
}

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

func (a *auth) Authenticate(email, password string) (*account_new.User, bool) {
	// FIXME remove para auth struct
	store, err := account_new.NewStorable()
	if err != nil {
		Logger.Warn(err.Error())
		return nil, false
	}
	defer store.Close()

	user, err := store.FindUserByEmail(email)
	if err != nil {
		Logger.Info("Failed trying to find the user '%s' to log in. Original Error: '%s'.", email, err.Error())
		return nil, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		Logger.Info("User '%s' is trying to log in with invalid password.", email)
		return nil, false
	}

	return &user, true
}

func (a *auth) CreateUserToken(user *account_new.User) (*ApiToken, error) {
	return createToken(user)
}

func (a *auth) UserFromToken(token string) (*account_new.User, error) {
	h := strings.Split(token, " ")
	if len(h) == 2 {
		apiToken := &ApiToken{Type: h[0], Token: h[1]}

		if apiToken.Type == TOKEN_TYPE {
			var user account_new.User

			err := decodeToken(apiToken, &user)
			if err != nil {
				return nil, err
			}
			if user.Email == "" {
				return nil, errors.ErrTokenNotFound
			}

			return &user, nil
		}
	}

	return nil, errors.ErrInvalidTokenFormat
}

func (a *auth) RevokeUserToken(token string) {
	user, err := a.UserFromToken(token)
	if err == nil {
		key := fmt.Sprintf("%s: %s", TOKEN_TYPE, user.Email)
		deleteToken(key)

		h := strings.Split(token, " ")
		deleteToken(h[1])
	}
}
