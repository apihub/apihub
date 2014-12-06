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
	ExpiresInSeconds  = 24 * 3600
	ExpiresTokenCache = ExpiresInSeconds - 10 // time in seconds to remove from expire time.
	TokenType         = "Token"
)

type Token interface {
	GetToken() (tokenType string, token string, error error)
	TokenFor() *TokenInfo
	GenerateToken() *TokenInfo
}

type TokenInfo struct {
	Token     string `json:"token"`
	Type      string `json:"token_type"`
	Expires   int    `json:"expires"`
	CreatedAt string `bson:"created_at" json:"created_at"`
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

// First, try to retrieve an existing token for the user. Return a new one if not found.
func TokenFor(user *account.User) *TokenInfo {
	token, err := getToken("token:" + user.Email)
	if err != nil {
		fmt.Print(err.Error())
	}
	var t TokenInfo
	redis.ScanStruct(token, &t)
	if t.Token == "" {
		conn, err := db.Conn()
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()
		t := GenerateToken(user)
		go conn.Tokens("token:"+user.Email, ExpiresTokenCache, structs.Map(t))
		return t
	}
	return &t
}

func GenerateToken(user *account.User) *TokenInfo {
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}

	token := &TokenInfo{
		Token:     base64.URLEncoding.EncodeToString(rb),
		Expires:   ExpiresInSeconds,
		CreatedAt: time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z07:00"),
		Type:      "Token",
	}
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	go conn.Tokens(token.Token, token.Expires, structs.Map(user))
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
