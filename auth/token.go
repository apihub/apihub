// Package auth provides an interface to authenticate user against the api.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
	. "github.com/backstage/backstage/errors"
	"github.com/fatih/structs"
	"github.com/garyburd/redigo/redis"
)

const (
	ExpiresInSeconds  = 24 * 3600
	ExpiresTokenCache = ExpiresInSeconds - 10 // time in seconds to remove from expire time.
	TokenType         = "Token"
)

type Token interface {
	GetUserFromToken(auth string) (user *account.User, error error)
	TokenFor(user *account.User) *TokenInfo
	RevokeTokensFor(user *account.User)
	GenerateToken(user *account.User) *TokenInfo
}

type TokenInfo struct {
	Token     string `json:"token"`
	Type      string `json:"token_type"`
	Expires   int    `json:"expires"`
	CreatedAt string `bson:"created_at" json:"created_at"`
}

// Convert a Token in a user.
// Given a token, find the user.
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
			if len(u) == 0 {
				return nil, ErrTokenNotFound
			}
			var user account.User
			if err := redis.ScanStruct(u, &user); err != nil {
				fmt.Print(err)
				return nil, err
			}
			//panic(u)
			return &user, nil
		}
	}

	return nil, ErrInvalidTokenFormat
}

// Return an auth token for the given user.
// This token should be used when calling the HTTP Api.
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

// Revoke the tokens for the given user.
// One example of use is when the user destroy its own account.
func RevokeTokensFor(user *account.User) {
	conn, err := db.Conn()
	defer conn.Close()
	ti := "token:" + user.Email
	token, err := getToken(ti)
	if err != nil {
		fmt.Print(err.Error())
	}
	var t TokenInfo
	redis.ScanStruct(token, &t)
	conn.DeleteToken(t.Token)
	conn.DeleteToken(ti)
}

// Generate a token for given user.
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
	conn.Tokens(token.Token, token.Expires, structs.Map(user))
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
