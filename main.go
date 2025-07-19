package main

import (
	"context"
)

func main() {
	router := NewBackendRouter()
	db := InitDB(DatabaseConnectionString)

	config := NewConfig("./config.yaml")

	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextAppName, "basic")

	router.Use(NewDBMiddleware(db))
	router.Use(NewCORSMiddleware())
	router.Use(NewLoggerMiddleware())
	router.UnprotectedHandle("GET /ping", HandlePing)

	router.UseProtected(NewUserIDMiddleware())
	router.ProtectedHandle("GET /profile", HandleProtected)

	webService := NewWebService(config, router)
	webService.Run(ctx)
}
