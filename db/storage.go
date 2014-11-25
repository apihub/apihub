package storage

import (
	"gopkg.in/mgo.v2"

	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/db/storage"
)

const (
	DefaultDatabaseHost = "127.0.0.1:27017"
	DefaultDatabaseName = "backstage"
)

type Storage struct {
	*storage.Storage
}

func conn() (*storage.Storage, error) {
	databaseHost, err := config.GetString("database:host")
	if err != nil {
		databaseHost = DefaultDatabaseHost
	}

	databaseName, err := config.GetString("database:name")
	if err != nil {
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
	subdomainIndex := mgo.Index{Key: []string{"subdomain"}, Unique: true}
	collection := storage.Collection("services")
	collection.EnsureIndex(subdomainIndex)
	return collection
}
