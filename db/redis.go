// Package db encapsulates connection with different storages: MongoDB, Redis.
package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/tsuru/config"
)

// const DefaultRedisHost = "apihub_redis_1:6379"
const DefaultRedisHost = "127.0.0.1:6379"

var redisPool *redis.Pool

type RedisClient struct {
	conn redis.Conn
	redis.PubSubConn
	sync.Mutex
}

func (c *RedisClient) Close() {
	c.conn.Close()
}

func getRedis() *redis.Pool {
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
		MaxActive:   4,
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

func NewRedisClient() *RedisClient {
	conn := getRedis().Get()
	client := &RedisClient{conn, redis.PubSubConn{conn}, sync.Mutex{}}
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			client.Lock()
			client.conn.Flush()
			client.Unlock()
		}
	}()
	return client
}

func DelCache(key string) (interface{}, error) {
	conn := NewRedisClient().conn
	defer conn.Close()
	result, err := conn.Do("DEL", key)
	if err != nil {
		fmt.Println("ERROR REDIS:", err)
		return nil, err
	}
	return result, nil
}

func GetHCache(key string) ([]interface{}, error) {
	conn := getRedis().Get()
	defer conn.Close()
	keyValue, err := conn.Do("HGETALL", key)
	if err != nil {
		fmt.Println("ERROR REDIS:", err)
		return nil, err
	}
	return keyValue.([]interface{}), nil
}

func AddHCache(key string, expires int, data map[string]interface{}) {
	conn := NewRedisClient().conn
	defer conn.Close()

	if _, err := conn.Do("HMSET", redis.Args{key}.AddFlat(data)...); err != nil {
		fmt.Print(err)
	}
	conn.Do("EXPIRE", key, expires)
}
