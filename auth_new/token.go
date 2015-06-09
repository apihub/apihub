package auth_new

import (
	"fmt"
	"time"

	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/db"
	. "github.com/backstage/backstage/log"
	"github.com/backstage/backstage/util"
	"github.com/fatih/structs"
)

const (
	EXPIRES_IN_SECONDS  = 24 * 3600
	EXPIRES_TOKEN_CACHE = 10 // time in seconds to remove from expire time.
	TOKEN_TYPE          = "Token"
)

type ApiToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	Expires   int    `json:"expires"`
	CreatedAt string `json:"created_at"`
}

// FIXME: need to improve this.
func createToken(user *account_new.User) (*ApiToken, error) {
	random := util.GenerateRandomStr(32)

	api := &ApiToken{
		Token:     random,
		Expires:   EXPIRES_IN_SECONDS,
		CreatedAt: time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z07:00"),
		Type:      TOKEN_TYPE,
	}

	key := fmt.Sprintf("%s: %s", api.Type, user.Email)

	db.Cache.Set(key, nil, time.Duration(api.Expires-EXPIRES_TOKEN_CACHE)*time.Minute)
	db.HMSET(key, api.Expires, structs.Map(api))

	db.Cache.Set(api.Token, nil, time.Duration(api.Expires))
	db.HMSET(api.Token, api.Expires, structs.Map(user))

	return api, nil
}

// FIXME: need to improve this.
func decodeToken(api *ApiToken, t interface{}) error {
	conn, err := db.Conn()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer conn.Close()

	return conn.GetTokenValue(api.Token, t)
}

// FIXME: need to improve this.
func deleteToken(key string) error {
	conn, err := db.Conn()
	if err != nil {
		Logger.Warn(err.Error())
		return err
	}
	defer conn.Close()

	conn.DeleteToken(key)
	return nil
}
