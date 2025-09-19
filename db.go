package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var _db *bun.DB

func GetDB() *bun.DB {
	return _db
}

func InitDB(dsn string) *bun.DB {
	sqldb, _ := sql.Open(sqliteshim.ShimName, dsn)
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(1 * time.Minute)
	sqldb.SetConnMaxIdleTime(1 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	_db = db

	return db
}

func createModels(ctx context.Context, db *bun.DB) error {
	if _, err := db.NewCreateTable().Model((*UserModel)(nil)).Exec(ctx); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}
	if _, err := db.NewCreateTable().Model((*SessionModel)(nil)).Exec(ctx); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}
	return nil
}

func DatabaseMiddleware(db *bun.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextDB, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
