package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/uptrace/bun"
)

func GetSession(r *http.Request) *sessions.Session {
	if sess, ok := r.Context().Value(ContextSession).(*sessions.Session); ok {
		return sess
	}
	return nil
}

func GetUser(r *http.Request) *UserModel {
	if u, ok := r.Context().Value(ContextUser).(*UserModel); ok {
		return u
	}
	panic(fmt.Errorf("user not found in context"))
}

func AuthMiddleware(db *bun.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := GetSession(r)
			if session == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_userID, err := uuid.Parse(userID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var user UserModel
			err = db.NewSelect().
				Model(&user).
				Where("id = ?", _userID).
				Scan(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUser, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
