package storage

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

var (
	pool = make(map[string]*Storage)
)

type Storage struct {
	session      *mgo.Session
	databaseName string
}

func open(address, databaseName string) (*Storage, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}

	storage := &Storage{session: session, databaseName: databaseName}
	poolKey := poolKey(address, databaseName)
	pool[poolKey] = storage
	return storage, nil
}

func Open(address, databaseName string) (*Storage, error) {
	poolKey := poolKey(address, databaseName)
	if storage, ok := pool[poolKey]; ok {
		return storage, nil
	}

	return open(address, databaseName)
}

func poolKey(address string, databaseName string) string {
	return fmt.Sprintf("%s:%s", address, databaseName)
}
