package mongore

import (
	"fmt"

	"github.com/tsuru/tsuru/db/storage"
)

const (
	DefaultDatabaseHost = "127.0.0.1:27017"
	DefaultDatabaseName = "apihub"
)

func (m *Mongore) openSession() *storage.Storage {
	if m.config.Host == "" {
		m.config.Host = DefaultDatabaseHost
	}
	if m.config.DatabaseName == "" {
		m.config.DatabaseName = DefaultDatabaseName
	}

	s, err := storage.Open(m.config.Host, m.config.DatabaseName)
	if err != nil {
		panic(fmt.Sprintf("Error while establishing connection to MongoDB: %s", err.Error()))
		return nil
	}
	return s
}
