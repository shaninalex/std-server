package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/shaninalex/std-server/pkg"
)

func DatabaseMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), pkg.ContextDB, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
