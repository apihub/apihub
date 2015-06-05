package mongore

import (
	"github.com/tsuru/tsuru/db/storage"
)

const (
	DefaultDatabaseHost = "127.0.0.1:27017"
	DefaultDatabaseName = "backstage"
)

func getConnection(config Config) (*storage.Storage, error) {
	if config.Host == "" {
		config.Host = DefaultDatabaseHost
	}
	if config.DatabaseName == "" {
		config.DatabaseName = DefaultDatabaseName
	}

	s, err := storage.Open(config.Host, config.DatabaseName)
	return s, err
}
