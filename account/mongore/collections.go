package mongore

import (
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

func (m *Mongore) Users() *storage.Collection {
	session := m.openSession()
	// defer session.Close()
	emailIndex := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	collection := session.Collection("users")
	collection.EnsureIndex(emailIndex)
	return collection
}

func (m *Mongore) Teams() *storage.Collection {
	session := m.openSession()
	// defer session.Close()
	index := mgo.Index{Key: []string{"alias"}, Unique: true, Background: false}
	collection := session.Collection("teams")
	collection.EnsureIndex(index)
	return collection
}

func (m *Mongore) Services() *storage.Collection {
	session := m.openSession()
	// defer session.Close()
	index := mgo.Index{Key: []string{"subdomain"}, Unique: true, Background: false}
	collection := session.Collection("services")
	collection.EnsureIndex(index)
	return collection
}
