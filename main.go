package main

import (
	"context"
)

func main() {
	router := NewRouter()
	config := NewConfig("./config.yaml")
	db := InitDB(config.String("db.dsn"))

	ctx := context.Background()
	if err := createModels(ctx, db); err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, ContextAppName, "basic")

	store := NewCookieStoreDatabase(db, []byte("very-secret-key"))

	router.Use(DatabaseMiddleware(db))
	router.Use(CorsMiddleware)
	router.Use(LoggerMiddleware)

	router.GET("/ping", HandlePing)
	router.POST("/login", HandleLogin)
	router.POST("/register", HandleRegister)

	router.Use(AuthMiddleware(db))
	router.Use(SessionMiddleware(store, "app_session"))
	router.GET("/profile", HandleProfile)

	webService := NewWebService(config, router)
	webService.Run(ctx)
}
