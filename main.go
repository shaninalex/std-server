package main

import (
	"context"
)

func main() {
	router := NewBackendRouter()
	config := NewConfig("./config.yaml")
	db := InitDB(config.String("db.dsn"))

	ctx := context.Background()
	if err := createModels(ctx, db); err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, ContextAppName, "basic")

	store := NewCookieStoreDatabase(db, []byte("very-secret-key"))

	router.Use(NewDBMiddleware(db))
	router.Use(NewCORSMiddleware())
	router.Use(NewLoggerMiddleware())

	router.Public().Get("/ping", HandlePing)
	router.Public().Post("/login", HandleLogin)
	router.Public().Post("/register", HandleRegister)

	router.Private().Use(NewSessionMiddleware(store, "app_session"))
	router.Private().Use(NewAuthMiddleware(db))
	router.Private().Get("/profile", HandleProfile)

	webService := NewWebService(config, router)
	webService.Run(ctx)
}
