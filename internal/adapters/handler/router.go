package handler

import (
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/api", getApiRouter())

	return r
}

func getApiRouter() *chi.Mux {
	api := chi.NewRouter()

	v1 := chi.NewRouter()

	registerHealthCheckRoute(v1)

	api.Mount("/v1", v1)

	return api
}

func registerHealthCheckRoute(r *chi.Mux) {
	r.Get("/health", healthCheckHandler)
}
