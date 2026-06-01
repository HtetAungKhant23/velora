package handler

import (
	"github.com/go-chi/chi"
)

type RouterDeps struct {
	AuthHandler *AuthHandler
}

func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/api", getApiRouter(deps))

	return r
}

func getApiRouter(deps RouterDeps) *chi.Mux {
	api := chi.NewRouter()

	v1 := chi.NewRouter()

	registerHealthCheckRoute(v1)
	registerAuthRoute(v1, deps.AuthHandler)

	api.Mount("/v1", v1)

	return api
}

func registerHealthCheckRoute(r *chi.Mux) {
	r.Get("/health", healthCheckHandler)
}

func registerAuthRoute(r *chi.Mux, h *AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})
}
