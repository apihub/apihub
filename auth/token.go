package auth

import (
	"errors"
	"strings"
)

func GetToken(auth string) (tokentype string, token string, error error) {
	var (
		tt string
		t  string
	)
	a := strings.Split(auth, " ")
	if len(a) == 2 {
		tt, t = a[0], a[1]
		return tt, t, nil
	}

	return tt, t, errors.New("Invalid token format.")
}
