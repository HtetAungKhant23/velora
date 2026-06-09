package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/HtetAungKhant23/velora/internal/adapters/handler"
	"github.com/HtetAungKhant23/velora/internal/adapters/handler/middleware"
	"github.com/HtetAungKhant23/velora/internal/adapters/processor"
	"github.com/HtetAungKhant23/velora/internal/adapters/repository"
	"github.com/HtetAungKhant23/velora/internal/adapters/storage"
	"github.com/HtetAungKhant23/velora/internal/adapters/token"
	"github.com/HtetAungKhant23/velora/internal/config"
	"github.com/HtetAungKhant23/velora/internal/core/services"
)

func main() {
	cfg := config.Load()

	db, err := repository.OpenDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect postgres:", "err", err, "dsn", cfg.DatabaseURL)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("postgres connected:", "dsn", cfg.DatabaseURL)

	// impl secondary adapters (outbound ports)
	userRepo := repository.NewUserRepository(db)
	imageRepo := repository.NewImageRepository(db)

	imageProcessor := processor.NewImageProcessor()
	imageStorage, err := storage.NewLocalStorage(cfg.StorageBaseDir, cfg.StorageBaseURL)
	if err != nil {
		slog.Error("failed to init image storage:", "err", err)
		os.Exit(1)
	}
	slog.Info("image storage ready:", "dir", cfg.StorageBaseDir)

	tokenSvc := token.NewJWTTokenService(cfg.JWTSecret, cfg.JWTExpiry)
	authGuard := middleware.NewAuthGuard(tokenSvc)

	// impl application service (inbound ports)
	userSvc := services.NewUserService(userRepo, tokenSvc)
	imageSvc := services.NewImageService(imageRepo, imageStorage, imageProcessor)

	// impl primary adapters (http handlers)
	authHandler := handler.NewAuthHandler(userSvc)
	imageHandler := handler.NewImageHandler(imageSvc)

	httpHandler := handler.NewRouter(handler.RouterDeps{
		AuthGuard:    authGuard,
		AuthHandler:  authHandler,
		ImageHandler: imageHandler,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Velora: image processing server listening on:", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server error:", "err", err)
		os.Exit(1)
	}
}
