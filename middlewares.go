package main

import (
	"context"
	"log"
	"net/http"
)

type IMiddleware interface {
	Wrap(http.Handler) http.Handler
}

type LoggerMiddleware struct{}

func (m *LoggerMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

type UserIDMiddleware struct{}

func (m *UserIDMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := "user_id__2091312039123"
		ctx := context.WithValue(r.Context(), ContextUserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewUserIDMiddleware() *UserIDMiddleware {
	return &UserIDMiddleware{}
}

type CORSMiddleware struct{}

func (m *CORSMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}
