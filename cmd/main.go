package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HtetAungKhant23/velora/internal/adapters/handler"
	"github.com/HtetAungKhant23/velora/internal/config"
)

func main() {
	cfg := config.Load()

	httpHandler := handler.NewRouter()

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
