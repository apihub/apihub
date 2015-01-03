package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
	. "github.com/backstage/backstage/log"
)

func Login(u *User) (*User, error) {
	user, err := FindUserByEmail(u.Email)
	if err != nil {
		Logger.Info("Failed trying to find the user '%s' to log in. Original Error: '%s'.", u.Email, err.Error())
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		Logger.Info("User '%s' trying to login with invalid password.", u.Email)
		return nil, err
	}
	return user, nil
}

func LoginAndGetToken(u *User) (*auth.TokenInfo, error) {
	user, err := Login(u)
	if err != nil {
		return nil, err
	}
	token := auth.TokenFor(user)
	return token, nil
}
