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
	"github.com/karlseguin/ccache"
)

const (
	ExpiresInSeconds  = 24 * 3600
	ExpiresTokenCache = ExpiresInSeconds - 10 // time in seconds to remove from expire time.
	TokenType         = "Token"
)

var Cache = ccache.New(ccache.Configure())

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
		tt       string
		tokenKey string
	)
	a := strings.Split(auth, " ")
	if len(a) == 2 {
		tt, tokenKey = a[0], a[1]
		if tt == TokenType {
			if item := Cache.Get(tokenKey); item != nil {
				return item.Value().(*account.User), nil
			}
			u, err := get(tokenKey)
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
	if item := Cache.Get(tokenKey); item != nil {
		return item.Value().(*TokenInfo)
	}
	token, err := get(tokenKey)
	if err != nil {
		fmt.Print(err.Error())
	}
	var t TokenInfo
	redis.ScanStruct(token, &t)
	if t.Token == "" {
		conn, err := db.Conn()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer conn.Close()
		t := GenerateToken(user)
		Cache.Set(tokenKey, t, time.Minute*10)
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
	token, err := get(ti)
	if err != nil {
		fmt.Print(err.Error())
	}
	var t TokenInfo
	redis.ScanStruct(token, &t)
	conn.DeleteToken(t.Token)
	Cache.Delete(t.Token)
	conn.DeleteToken(ti)
	Cache.Delete(ti)
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
	Cache.Set(token.Token, user, time.Minute*10)
	conn.Tokens(token.Token, token.Expires, structs.Map(user))
	return token
}

func get(token string) ([]interface{}, error) {
	conn, err := db.Conn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	return conn.GetTokenValue(token)
}
