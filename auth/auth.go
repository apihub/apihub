package auth

import (
	"fmt"
	"strings"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	. "github.com/backstage/backstage/log"
	"github.com/backstage/backstage/util"
)

const (
	EXPIRES_IN_SECONDS = 24 * 3600
	TOKEN_TYPE         = "Token"
)

type Authenticatable interface {
	Authenticate(email, password string) (*account.User, bool)
	CreateUserToken(*account.User) (*account.TokenInfo, error)
	UserFromToken(token string) (*account.User, error)
	RevokeUserToken(token string) error
}

type auth struct {
	store account.Storable
}

func NewAuth(store account.Storable) *auth {
	return &auth{store: store}
}

func (a *auth) Authenticate(email, password string) (*account.User, bool) {
	user, err := a.store.FindUserByEmail(email)
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

func (a *auth) CreateUserToken(user *account.User) (*account.TokenInfo, error) {
	api := account.TokenInfo{
		CreatedAt: time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z07:00"),
		Expires:   EXPIRES_IN_SECONDS,
		Type:      TOKEN_TYPE,
		Token:     util.GenerateRandomStr(32),
		User:      user,
	}

	err := a.store.CreateToken(api)
	if err != nil {
		return nil, err
	}

	return &api, err
}

func (a *auth) UserFromToken(token string) (*account.User, error) {
	h := strings.Split(token, " ")
	if len(h) == 2 {
		apiToken := account.TokenInfo{Type: h[0], Token: h[1]}

		if apiToken.Type == TOKEN_TYPE {
			var user account.User

			if err := a.store.DecodeToken(apiToken.Token, &user); err != nil {
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

func (a *auth) RevokeUserToken(token string) error {
	user, err := a.UserFromToken(token)

	if err == nil && user.Email != "" {
		key := fmt.Sprintf("%s: %s", TOKEN_TYPE, user.Email)
		err = a.store.DeleteToken(key)
		if err == nil {
			h := strings.Split(token, " ")
			err = a.store.DeleteToken(h[1])
		}
	}

	return err
}
