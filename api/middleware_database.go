package api

import (
	"context"
	"net/http"

	"github.com/shaninalex/std-server/pkg"
	"github.com/uptrace/bun"
)

func DatabaseMiddleware(db *bun.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), pkg.ContextDB, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
