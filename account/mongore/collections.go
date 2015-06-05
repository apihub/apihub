package mongore

import (
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

func (m *Mongore) Users() *storage.Collection {
	emailIndex := mgo.Index{Key: []string{"email"}, Unique: true, Background: false}
	usernameIndex := mgo.Index{Key: []string{"username"}, Unique: true, Background: false}
	collection := m.store.Collection("users")
	collection.EnsureIndex(emailIndex)
	collection.EnsureIndex(usernameIndex)
	return collection
}
