package main

import (
	"context"
	"fmt"
	"net/http"
)

type SessionMiddleware struct {
	store       *CookieStoreDatabase
	sessionName string
}

func NewSessionMiddleware(store *CookieStoreDatabase, sessionName string) *SessionMiddleware {
	return &SessionMiddleware{
		store:       store,
		sessionName: sessionName,
	}
}

func (s *SessionMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, s.sessionName)
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), ContextSession, session)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		if session.IsNew || len(session.Values) > 0 {
			if err := s.store.Save(r, w, session); err != nil {
				fmt.Println("Failed to save session:", err)
			}
		}
	})
}
