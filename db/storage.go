package db

import (
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

const (
	DefaultDatabaseHost = "127.0.0.1:27017"
	DefaultDatabaseName = "backstage"
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

func (storage *Storage) Users() *storage.Collection {
	emailIndex := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	collection := storage.Collection("users")
	collection.EnsureIndex(emailIndex)
	return collection
}

func (storage *Storage) Teams() *storage.Collection {
	aliasIndex := mgo.Index{Key: []string{"alias"}, Unique: true, Background: false}
	collection := storage.Collection("teams")
	collection.EnsureIndex(aliasIndex)
	return collection
}

func (storage *Storage) Tokens(key string, expires int, data map[string]interface{}) {
	addHCache(key, expires, data)
}

func (storage *Storage) GetTokenValue(key string) ([]interface{}, error) {
	return getHCache(key)
}

func (storage *Storage) DeleteToken(key string) (interface{}, error) {
	return delCache(key)
}
