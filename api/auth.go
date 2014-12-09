package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
)

func Login(email string, password string) (*auth.TokenInfo, error) {
	var user *account.User
	user, err := account.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	token := auth.TokenFor(user)
	return token, err
}
