package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const (
	ExpiresInSeconds = 24 * 3600
)

type Token interface {
	GetToken() (tokenType string, token string, error error)
	GenerateToken() *TokenInfo
}

type TokenInfo struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

func GetToken(auth string) (tokenType string, token string, error error) {
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

func GenerateToken() *TokenInfo {
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}

	token := &TokenInfo{Token: base64.URLEncoding.EncodeToString(rb), Expires: ExpiresInSeconds}
	return token
}
