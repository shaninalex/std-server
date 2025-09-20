package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shaninalex/std-server/pkg"
)

func SessionMiddleware(store *pkg.CookieStoreDatabase, sessionName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, sessionName)
			if err != nil {
				http.Error(w, "Session error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), pkg.ContextSession, session)
			next.ServeHTTP(w, r.WithContext(ctx))
			if session.IsNew || len(session.Values) > 0 {
				if err := store.Save(r, w, session); err != nil {
					fmt.Println("Failed to save session:", err)
				}
			}
		})
	}
}
