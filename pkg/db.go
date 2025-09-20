package pkg

import (
	"database/sql"
	"time"
)

// TODO: remove global variable
var _db *sql.DB

func GetDB() *sql.DB {
	return _db
}

func InitDB(dsn string) *sql.DB {
	sqldb, _ := sql.Open("sqlite3", dsn)
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(1 * time.Minute)
	sqldb.SetConnMaxIdleTime(1 * time.Minute)

	_db = sqldb

	return sqldb
}
