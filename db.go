package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(1 * time.Minute)
	sqldb.SetConnMaxIdleTime(1 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	return db
}

type DBMiddleware struct {
	db *bun.DB
}

func NewDBMiddleware(db *bun.DB) *DBMiddleware {
	return &DBMiddleware{db: db}
}

func (m *DBMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextDB, m.db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DBFromContext(ctx context.Context) *bun.DB {
	db := ctx.Value(ContextDB).(*bun.DB)
	if db == nil {
		panic("postgres context is not set")
	}
	return db
}

func DBConnect(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(sqldb, pgdialect.New())
}
