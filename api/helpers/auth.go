package helpers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/albertoleal/backstage/account"
)

func SignIn(username string, password string) (*account.User, error) {
	var user *account.User
	user, err := account.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, err
}
