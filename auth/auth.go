package auth

import (
	"strings"
	"time"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
	"github.com/backstage/backstage/util"
)

type AuthenticationToken struct {
	storage account.Storable
}

func NewAuthenticationToken(storage account.Storable) *AuthenticationToken {
	return &AuthenticationToken{
		storage: storage,
	}
}

// Convert a Token in a user.
// Given a token, find the user.
func (at *AuthenticationToken) GetUserFromToken(auth string) (*account.User, error) {
	var (
		tt       string
		tokenKey string
	)

	a := strings.Split(auth, " ")
	if len(a) == 2 {
		tt, tokenKey = a[0], a[1]
		if tt == TokenType {
			u, err := at.storage.GetToken(account.TokenKey{Name: tokenKey})
			if err != nil {
				return nil, err
			}
			return u.(*account.User), nil
		}
	}

	return nil, errors.ErrInvalidTokenFormat
}

// Return an auth token for the given user.
// This token should be used when calling the HTTP Api.
// First, try to retrieve an existing token for the user. Return a new one if not found.
func (at *AuthenticationToken) TokenFor(user *account.User) (*account.TokenInfo, error) {
	tokenKey := account.TokenKey{Name: "token:" + user.Email}
	ti, err := at.storage.GetToken(tokenKey)
	if _, ok := err.(*errors.NotFoundError); ok {
		t, err := at.generateToken()
		if err != nil {
			return &account.TokenInfo{}, err
		}
		err = at.storage.SaveToken(tokenKey, t)
		if err != nil {
			return &account.TokenInfo{}, err
		}
		err = at.storage.SaveToken(account.TokenKey{Name: t.Token}, user)
		if err != nil {
			return &account.TokenInfo{}, err
		}
		return t, nil
	}
	tokenInfo := ti.(*account.TokenInfo)
	return tokenInfo, nil
}

// Revoke the tokens for the given user.
// One example of use is when the user destroy its own account.
func (at *AuthenticationToken) RevokeTokenFor(user *account.User) {
	tokenKey := account.TokenKey{Name: "token:" + user.Email}
	ti, err := at.storage.GetToken(tokenKey)
	if _, ok := err.(*errors.NotFoundError); ok {
		return
	}
	tokenInfo := ti.(*account.TokenInfo)
	at.storage.DeleteToken(account.TokenKey{Name: tokenInfo.Token})
	at.storage.DeleteToken(tokenKey)
}

// Generate a token for given user.
func (at *AuthenticationToken) generateToken() (*account.TokenInfo, error) {
	tok := util.GenerateRandomStr(32)
	token := &account.TokenInfo{
		Token:     tok,
		Expires:   ExpiresInSeconds,
		CreatedAt: time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z07:00"),
		Type:      "Token",
	}
	return token, nil
}
