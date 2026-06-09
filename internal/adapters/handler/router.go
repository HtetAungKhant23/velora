package handler

import (
	"net/http"

	"github.com/HtetAungKhant23/velora/internal/adapters/handler/middleware"
	"github.com/go-chi/chi"
)

type RouterDeps struct {
	AuthGuard    *middleware.AuthGuard
	AuthHandler  *AuthHandler
	ImageHandler *ImageHandler
	StaticDir    string
}

func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	registerStaticRoute(r, deps.StaticDir)

	r.Mount("/api", getApiRouter(deps))

	return r
}

func getApiRouter(deps RouterDeps) *chi.Mux {
	api := chi.NewRouter()

	v1 := chi.NewRouter()

	registerHealthCheckRoute(v1)
	registerAuthRoute(v1, deps.AuthHandler, deps.AuthGuard)
	registerImageRoute(v1, deps.ImageHandler, deps.AuthGuard)

	api.Mount("/v1", v1)

	return api
}

func registerHealthCheckRoute(r chi.Router) {
	r.Get("/health", healthCheckHandler)
}

func registerAuthRoute(r chi.Router, h *AuthHandler, authGuard *middleware.AuthGuard) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(authGuard.Verify)

			r.Get("/me", h.Me)
		})
	})
}

func registerImageRoute(r chi.Router, h *ImageHandler, authGuard *middleware.AuthGuard) {
	r.Route("/images", func(r chi.Router) {
		r.Use(authGuard.Verify)

		r.Post("/upload", h.Upload)
	})
}

func registerStaticRoute(r chi.Router, baseDir string) {
	if baseDir == "" {
		return
	}

	fileSvcHandler := http.FileServer(http.Dir(baseDir))
	r.Get("/file/*", http.StripPrefix("/file/", fileSvcHandler).ServeHTTP)
}
