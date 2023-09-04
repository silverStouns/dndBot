package database

import (
	"database/sql"
	"dndBot/internal/pkg/logger"
	_ "modernc.org/sqlite"
	"sync"
)

type DBConnector struct {
	Connector            *sql.DB
	requestsSourcesCache map[string]string
	cacheMutexLocker     sync.Mutex
}

func GetNewDBConnector() (*DBConnector, error) {
	connector, err := sql.Open("sqlite", "botDatabase.db")
	if err != nil {
		return nil, err
	}
	if err = connector.Ping(); err != nil {
		return nil, err
	}
	connector.SetMaxIdleConns(5)
	logger.Trace("Database was connected")

	return &DBConnector{
		Connector:            connector,
		requestsSourcesCache: make(map[string]string),
		cacheMutexLocker:     sync.Mutex{},
	}, err
}
