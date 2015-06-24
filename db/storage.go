package db

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

const (
	DefaultDatabaseHost = "127.0.0.1:27017"
	DefaultDatabaseName = "backstage_development"
)

type Storage struct {
	*storage.Storage
}

func conn() (*storage.Storage, error) {
	databaseHost, _ := config.GetString("database:host")
	if databaseHost == "" {
		databaseHost = DefaultDatabaseHost
	}

	databaseName, _ := config.GetString("database:name")
	if databaseName == "" {
		databaseName = DefaultDatabaseName
	}

	return storage.Open(databaseHost, databaseName)
}

func Conn() (*Storage, error) {
	var (
		strg Storage
		err  error
	)

	strg.Storage, err = conn()
	return &strg, err
}

func (storage *Storage) Services() *storage.Collection {
	collection := storage.Collection("services")
	return collection
}

func (storage *Storage) Plugins() *storage.Collection {
	index := mgo.Index{Key: []string{"service", "name"}, Unique: true, Background: false}
	collection := storage.Collection("plugins_config")
	collection.EnsureIndex(index)
	return collection
}

func (storage *Storage) Users() *storage.Collection {
	emailIndex := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	usernameIndex := mgo.Index{Key: []string{"username"}, Unique: true, Background: false}
	collection := storage.Collection("users")
	collection.EnsureIndex(emailIndex)
	collection.EnsureIndex(usernameIndex)
	return collection
}

func (storage *Storage) Teams() *storage.Collection {
	aliasIndex := mgo.Index{Key: []string{"alias"}, Unique: true, Background: false}
	collection := storage.Collection("teams")
	collection.EnsureIndex(aliasIndex)
	return collection
}

func (storage *Storage) Clients() *storage.Collection {
	idIndex := mgo.Index{Key: []string{"id"}, Unique: true, Background: false}
	collection := storage.Collection("clients")
	collection.EnsureIndex(idIndex)
	return collection
}

func (storage *Storage) Authorizations() *storage.Collection {
	codeIndex := mgo.Index{Key: []string{"code"}}
	collection := storage.Collection("authorizations")
	collection.EnsureIndex(codeIndex)
	return collection
}

func (storage *Storage) Accesses() *storage.Collection {
	accessIndex := mgo.Index{Key: []string{"accesstoken"}}
	redirectIndex := mgo.Index{Key: []string{"redirecturi"}}
	collection := storage.Collection("accesses")
	collection.EnsureIndex(accessIndex)
	collection.EnsureIndex(redirectIndex)
	return collection
}

func (storage *Storage) Tokens(key string, expires int, data map[string]interface{}) {
	Cache.Set(key, nil, time.Duration(expires-10)*time.Minute)
	addHCache(key, expires, data)
}

func (storage *Storage) GetTokenValue(key string, t interface{}) error {
	var (
		data []interface{}
		err  error
	)

	if item := Cache.Get(key); item != nil && item.Value() != nil {
		if !item.Expired() {
			data = item.Value().([]interface{})
		}
	}
	if len(data) == 0 {
		data, err = getHCache(key)
		if err != nil {
			return err
		}
	}

	if err = redis.ScanStruct(data, t); err != nil {
		fmt.Print(err)
		return err
	}
	if len(data) > 0 {
		Cache.Replace(key, data)
	}
	return nil
}

func (storage *Storage) DeleteToken(key string) (interface{}, error) {
	Cache.Delete(key)
	return delCache(key)
}
