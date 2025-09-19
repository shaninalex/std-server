package main

import (
	"context"
	"fmt"
	"net/http"
)

func SessionMiddleware(store *CookieStoreDatabase, sessionName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, sessionName)
			if err != nil {
				http.Error(w, "Session error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), ContextSession, session)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

			if session.IsNew || len(session.Values) > 0 {
				if err := store.Save(r, w, session); err != nil {
					fmt.Println("Failed to save session:", err)
				}
			}
		})
	}
}
