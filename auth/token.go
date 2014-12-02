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
	"github.com/fatih/structs"
	"github.com/garyburd/redigo/redis"
)

const (
	ExpiresInSeconds = 24 * 3600
	TokenType        = "Token"
)

type Token interface {
	GetToken() (tokenType string, token string, error error)
	GenerateToken() *TokenInfo
}

type TokenInfo struct {
	User      *account.User `bson:"user" json:"-"`
	Token     string        `json:"token"`
	Type      string        `json:"token_type"`
	Expires   int           `json:"expires"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

func GetUserFromToken(auth string) (user *account.User, error error) {
	var (
		tt string
		t  string
	)
	a := strings.Split(auth, " ")
	if len(a) == 2 {
		tt, t = a[0], a[1]
		if tt == TokenType {
			u, err := getToken(t)
			if err != nil {
				return nil, err
			}
			var user account.User
			if err := redis.ScanStruct(u, &user); err != nil {
				fmt.Print(err)
			}
			return &user, nil
		}
	}

	return nil, errors.New("Invalid token format.")
}

func GenerateToken(user *account.User) *TokenInfo {
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}

	token := &TokenInfo{User: user, Token: base64.URLEncoding.EncodeToString(rb),
		Expires: ExpiresInSeconds, Type: "Token", CreatedAt: time.Now()}
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	conn.Tokens(token.Token, token.Expires, structs.Map(token.User))
	return token
}

func getToken(token string) ([]interface{}, error) {
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	return conn.GetTokenValue(token)
}
