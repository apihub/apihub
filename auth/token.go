// Package auth provides an interface to authenticate user against the api.
package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
	. "github.com/backstage/backstage/errors"
	"github.com/backstage/backstage/util"
	"github.com/fatih/structs"
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
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	Expires   int    `json:"expires"`
	CreatedAt string `bson:"created_at" json:"created_at"`
}

//Return a representation of token but without sensitive data.
func (t *TokenInfo) ToString() string {
	token, _ := json.Marshal(t)
	return string(token)
}

// Convert a Token in a user.
// Given a token, find the user.
func GetUserFromToken(auth string) (user *account.User, error error) {
	var (
		tt       string
		tokenKey string
	)
	a := strings.Split(auth, " ")
	if len(a) == 2 {
		tt, tokenKey = a[0], a[1]
		if tt == TokenType {
			var user account.User
			err := get(tokenKey, &user)
			if err != nil {
				return nil, err
			}
			if user.Email == "" {
				return nil, ErrTokenNotFound
			}
			return &user, nil
		}
	}
	return nil, ErrInvalidTokenFormat
}

// Return an auth token for the given user.
// This token should be used when calling the HTTP Api.
// First, try to retrieve an existing token for the user. Return a new one if not found.
func TokenFor(user *account.User) *TokenInfo {
	tokenKey := "token:" + user.Email
	var t TokenInfo
	err := get(tokenKey, &t)
	if err != nil {
		fmt.Print(err.Error())
	}
	if t.Token == "" {
		conn, err := db.Conn()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer conn.Close()
		t := GenerateToken(user)
		go conn.Tokens(tokenKey, ExpiresTokenCache, structs.Map(t))
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
	var t TokenInfo
	err = get(ti, &t)
	if err != nil {
		fmt.Print(err.Error())
	}
	conn.DeleteToken(t.Token)
	conn.DeleteToken(ti)
}

// Generate a token for given user.
func GenerateToken(user *account.User) *TokenInfo {
	tok := util.GenerateRandomStr(32)
	token := &TokenInfo{
		Token:     tok,
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

func get(token string, t interface{}) error {
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	return conn.GetTokenValue(token, t)
}
