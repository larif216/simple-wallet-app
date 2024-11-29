package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func newDatabase(cfg DatabaseConfig) *sql.DB {
	var sqlDB *sql.DB
	var err error

	sqlDB, err = sql.Open(cfg.Driver, cfg.RWDataSourceName())

	if err != nil {
		panic(err)
	}

	return sqlDB
}
