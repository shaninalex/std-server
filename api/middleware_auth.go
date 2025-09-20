package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/shaninalex/std-server/pkg"
)

func GetSession(r *http.Request) *sessions.Session {
	if sess, ok := r.Context().Value(pkg.ContextSession).(*sessions.Session); ok {
		return sess
	}
	return nil
}

func GetUser(r *http.Request) *pkg.UserModel {
	if u, ok := r.Context().Value(pkg.ContextUser).(*pkg.UserModel); ok {
		return u
	}
	panic(fmt.Errorf("user not found in context"))
}

func AuthMiddleware(db *sql.DB) func(next http.Handler) http.Handler {
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

			user, err := pkg.GetUserByField(r.Context(), db, "id", _userID.String())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), pkg.ContextUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
