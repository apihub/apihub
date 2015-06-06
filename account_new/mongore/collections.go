package mongore

import (
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

func (m *Mongore) Users() *storage.Collection {
	emailIndex := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	collection := m.store.Collection("users")
	collection.EnsureIndex(emailIndex)
	return collection
}

func (m *Mongore) Teams() *storage.Collection {
	aliasIndex := mgo.Index{Key: []string{"alias"}, Unique: true, Background: false}
	collection := m.store.Collection("teams")
	collection.EnsureIndex(aliasIndex)
	return collection
}
