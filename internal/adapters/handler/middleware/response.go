package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/HtetAungKhant23/velora/internal/core/domain/shared"
)

type ErrResponse struct {
	Error string `json:"error"`
}

type envelope struct {
	Data any `json:"data"`
}

var (
	SuccessResponse = struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}
)

func ErrResp(msg string) ErrResponse {
	return ErrResponse{Error: msg}
}

func ReadJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}

func WriteOK(w http.ResponseWriter, v any) {
	WriteJSON(w, http.StatusOK, envelope{Data: v})
}

func WriteCreated(w http.ResponseWriter, v any) {
	WriteJSON(w, http.StatusCreated, envelope{Data: v})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrResponse{Error: message})
}

func MapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, ErrResp("resource not found"))
	case errors.Is(err, shared.ErrUnauthorized):
		WriteJSON(w, http.StatusUnauthorized, ErrResp("invalid credentials"))
	case errors.Is(err, shared.ErrForbidden):
		WriteJSON(w, http.StatusForbidden, ErrResp("access denied"))
	case errors.Is(err, shared.ErrAlreadyExist):
		WriteJSON(w, http.StatusConflict, ErrResp(err.Error()))
	case errors.Is(err, shared.ErrInvalidInput):
		WriteJSON(w, http.StatusBadRequest, ErrResp(err.Error()))
	default:
		slog.Error("unhandled domain error", "err", err)
		WriteJSON(w, http.StatusInternalServerError, ErrResp(fmt.Sprintf("internal server error: %s", err.Error())))
	}
}
