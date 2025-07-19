package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	serverAddr = "127.0.0.1:8000"
)

func main() {
	router := NewBackendRouter()
	db := InitDB(DatabaseConnectionString)

	router.Use(NewDBMiddleware(db))
	router.Use(NewCORSMiddleware())
	router.Use(NewLoggerMiddleware())
	router.UnprotectedHandle("GET /ping", HandlePing)

	router.UseProtected(NewUserIDMiddleware())
	router.ProtectedHandle("GET /profile", HandleProtected)

	srv := &http.Server{
		Handler:      router,
		Addr:         serverAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Run server on: ", serverAddr)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
