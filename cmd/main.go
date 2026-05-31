package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HtetAungKhant23/velora/internal/adapters/handler"
	"github.com/HtetAungKhant23/velora/internal/adapters/repository"
	"github.com/HtetAungKhant23/velora/internal/adapters/token"
	"github.com/HtetAungKhant23/velora/internal/config"
	"github.com/HtetAungKhant23/velora/internal/core/services"
)

func main() {
	cfg := config.Load()

	db, err := repository.OpenDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v and dsn: %s", err, cfg.DatabaseURL)
	}
	defer db.Close()
	log.Printf("postgres connected: %s", cfg.DatabaseURL)

	tokenSvc := token.NewJWTTokenService(cfg.JWTSecret, cfg.JWTExpiry)

	userRepo := repository.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, tokenSvc)
	authHandler := handler.NewAuthHandler(userSvc)

	httpHandler := handler.NewRouter(handler.RouterDeps{
		AuthHandler: authHandler,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Velora: image processing server listening on %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
