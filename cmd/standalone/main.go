package main

import (
	"context"

	"github.com/shaninalex/std-server/api"
	"github.com/shaninalex/std-server/pkg"
)

func main() {
	router := api.NewRouter()
	pkg.ApplyMigrationsEmbed("file:myapp.db?cache=shared&mode=rwc")
	db := pkg.InitDB("file:myapp.db?cache=shared&mode=rwc")

	ctx := context.Background()
	ctx = context.WithValue(ctx, pkg.ContextAppName, "basic")

	store := pkg.NewCookieStoreDatabase(db, []byte("very-secret-key"))

	router.Use(api.DatabaseMiddleware(db))
	router.Use(api.CorsMiddleware)
	router.Use(api.LoggerMiddleware)

	router.GET("/ping", api.HandlePing)
	router.POST("/login", api.HandleLogin)
	router.POST("/register", api.HandleRegister)

	router.Use(api.SessionMiddleware(store, "app_session"))
	router.Use(api.AuthMiddleware(db))
	router.GET("/profile", api.HandleProfile)

	webService := api.NewWebService(router)
	webService.Run(ctx)
}
