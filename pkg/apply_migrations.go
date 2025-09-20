package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(dsn string) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3", driver,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
