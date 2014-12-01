package db

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
)

const (
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
	netloc := "localhost:6379"
	password := ""
	pool := &redis.Pool{
		MaxActive:   24,
		MaxIdle:     12,
		IdleTimeout: 60 * time.Second,
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

func (storage *Storage) Tokens(keys map[string]string, expires int) {
	conn := GetRedis()
	defer conn.Close()
	conn.Send("MULTI")
	for key, value := range keys {
		conn.Send("SETEX", key, expires, value)
	}
	_, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
