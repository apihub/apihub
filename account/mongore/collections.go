package mongore

import (
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

type Storage struct {
	*storage.Storage
}

func (strg *Storage) Users() *storage.Collection {
	emailIndex := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	collection := strg.Collection("users")
	collection.EnsureIndex(emailIndex)
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
