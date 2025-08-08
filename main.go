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

	router.Public().Get("/ping", HandlePing)

	router.Private().Use(NewUserIDMiddleware())
	router.Private().Get("/profile", HandleProtected)

	webService := NewWebService(config, router)
	webService.Run(ctx)
}
