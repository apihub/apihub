package mongore

import (
	"fmt"

	"github.com/apihub/apihub/db"
	"github.com/garyburd/redigo/redis"
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

func (strg *Storage) Plugins() *storage.Collection {
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

func (strg *Storage) GetTokenValue(key string, t interface{}) error {
	var (
		data []interface{}
		err  error
	)

	if item := db.Cache.Get(key); item != nil && item.Value() != nil {
		if !item.Expired() {
			data = item.Value().([]interface{})
		}
	}
	if len(data) == 0 {
		data, err = db.GetHCache(key)
		if err != nil {
			return err
		}
	}

	if err = redis.ScanStruct(data, t); err != nil {
		fmt.Print(err)
		return err
	}
	if len(data) > 0 {
		db.Cache.Replace(key, data)
	}
	return nil
}

func (strg *Storage) DeleteToken(key string) (interface{}, error) {
	db.Cache.Delete(key)
	return db.DelCache(key)
}
