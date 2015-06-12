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
	index := mgo.Index{Key: []string{"alias"}, Unique: true, Background: false}
	collection := m.store.Collection("teams")
	collection.EnsureIndex(index)
	return collection
}

func (m *Mongore) Services() *storage.Collection {
	index := mgo.Index{Key: []string{"subdomain"}, Unique: true, Background: false}
	collection := m.store.Collection("services")
	collection.EnsureIndex(index)
	return collection
}
