package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
)

func Login(u *User) (*auth.TokenInfo, error) {
	user, err := FindUserByEmail(u.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return nil, err
	}
	token := auth.TokenFor(user)
	return token, nil
}
