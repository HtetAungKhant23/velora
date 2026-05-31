package handler

import (
	"net/http"

	"github.com/HtetAungKhant23/velora/internal/adapters/handler/middleware"
	"github.com/HtetAungKhant23/velora/internal/core/ports"
)

type AuthHandler struct {
	useCase ports.UserUseCase
}

func NewAuthHandler(uc ports.UserUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: uc,
	}
}

type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	payload := &RegisterPayload{}

	if err := middleware.ReadJSON(r, payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.useCase.Register(r.Context(), ports.RegisterCommand{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		middleware.MapDomainError(w, err)
		return
	}

	middleware.WriteCreated(w, result)
}
