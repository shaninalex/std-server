package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type WebService struct {
	config IConfig
	router IBackendRouter
}

func NewWebService(config IConfig, router IBackendRouter) *WebService {
	return &WebService{
		config: config,
		router: router,
	}
}

func (s *WebService) Run(ctx context.Context) {
	appName := ctx.Value(ContextAppName).(string)
	addr := fmt.Sprintf(":%d", s.config.Int(appName+".port"))

	srv := &http.Server{
		Handler:      s.router,
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Run Server on %s ...\n", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("Shutdown Server %s ...\n", addr)

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}
	log.Println("Server exiting")
}
