package mongore

import (
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

type Storage struct {
	*storage.Storage
}

func (strg *Storage) Apps() *storage.Collection {
	index := mgo.Index{Key: []string{"clientid"}, Unique: true, Background: false}
	collection := strg.Collection("apps")
	collection.EnsureIndex(index)
	return collection
}

func (strg *Storage) Users() *storage.Collection {
	index := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	collection := strg.Collection("users")
	collection.EnsureIndex(index)
	return collection
}

func (strg *Storage) Teams() *storage.Collection {
	index := mgo.Index{Key: []string{"alias"}, Unique: true, Background: false}
	collection := strg.Collection("teams")
	collection.EnsureIndex(index)
	return collection
}

func (strg *Storage) Services() *storage.Collection {
	index := mgo.Index{Key: []string{"subdomain"}, Unique: true, Background: false}
	collection := strg.Collection("services")
	collection.EnsureIndex(index)
	return collection
}

func (strg *Storage) PluginsConfig() *storage.Collection {
	index := mgo.Index{Key: []string{"service", "name"}, Unique: true, Background: false}
	collection := strg.Collection("plugins_config")
	collection.EnsureIndex(index)
	return collection
}

func (strg *Storage) Hooks() *storage.Collection {
	index := mgo.Index{Key: []string{"name", "team"}, Unique: true, Background: false}
	collection := strg.Collection("hooks")
	collection.EnsureIndex(index)
	return collection
}
