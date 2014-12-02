package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

const (
	DefaultRedisHost = "127.0.0.1:6379"
	DefaultDatabaseHost = "127.0.0.1:27017"
	DefaultDatabaseName = "backstage"
)

var (
	redisPool *redis.Pool
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

func GetRedisPool() *redis.Pool {
	if redisPool != nil {
		return redisPool
	}
	netloc, _ := config.GetString("redis:host")

	if netloc == "" {
		netloc = DefaultRedisHost
	}

	password, _ := config.GetString("redis:password")
	redisNumber, _ := config.GetInt("redis:number")

	pool := &redis.Pool{
		MaxActive:   2,
		MaxIdle:     2,
		IdleTimeout: 0,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", netloc)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := conn.Do("AUTH", password); err != nil {
					conn.Close()
					return nil, err
				}
			}

			if redisNumber > 0 {
				if _, err := conn.Do("SELECT", redisNumber); err != nil {
					conn.Close()
					return nil, err
				}
			}

			return conn, nil
		},
	}
	redisPool = pool
	return redisPool
}

func GetRedis() redis.Conn {
	pool := GetRedisPool()
	return pool.Get()
}

func (storage *Storage) Services() *storage.Collection {
	subdomainIndex := mgo.Index{Key: []string{"subdomain"}, Unique: true}
	collection := storage.Collection("services")
	collection.EnsureIndex(subdomainIndex)
	return collection
}

func (storage *Storage) Users() *storage.Collection {
	usernameIndex := mgo.Index{Key: []string{"username"}, Unique: true}
	collection := storage.Collection("users")
	collection.EnsureIndex(usernameIndex)
	return collection
}

func (storage *Storage) Groups() *storage.Collection {
	nameIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	collection := storage.Collection("groups")
	collection.EnsureIndex(nameIndex)
	return collection
}

func (storage *Storage) Tokens(key string, expires int, data map[string]interface{}) {
	addHCache(key, expires, data)
}

func (storage *Storage) GetTokenValue(key string) ([]interface{}, error) {
	return getHCache(key)
}

func getHCache(key string) ([]interface{}, error) {
	conn := GetRedis()
	defer conn.Close()
	keyValue, err := conn.Do("HGETALL", key)
	if err != nil {
		fmt.Println("ERROR:", err)
		return nil, err
	}
	return keyValue.([]interface{}), nil
}

func addHCache(key string, expires int, data map[string]interface{}) {
	conn := GetRedis()
	defer conn.Close()

	if _, err := conn.Do("HMSET", redis.Args{key}.AddFlat(data)...); err != nil {
		fmt.Print(err)
	}
	conn.Do("EXPIRE", key, expires)
}
