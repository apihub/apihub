package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/albertoleal/backstage/account"
	"github.com/albertoleal/backstage/db"
)

const (
	ExpiresInSeconds = 24 * 3600
)

type Token interface {
	GetToken() (tokenType string, token string, error error)
	GenerateToken() *TokenInfo
}

type TokenInfo struct {
	User      string    `bson:"username" json:"username"`
	Token     string    `json:"token"`
	Expires   int       `json:"expires"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
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

func GenerateToken(user *account.User) *TokenInfo {
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}

	token := &TokenInfo{User: user.Username, Token: base64.URLEncoding.EncodeToString(rb),
		Expires: ExpiresInSeconds, CreatedAt: time.Now()}
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	conn.Tokens(map[string]string{token.Token: token.User}, token.Expires)
	return token
}
