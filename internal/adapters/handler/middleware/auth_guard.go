package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/HtetAungKhant23/velora/internal/core/ports"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyEmail  contextKey = "email"
)

type AuthGuard struct {
	tokenService ports.TokenService
}

func NewAuthGuard(tokenSvc ports.TokenService) *AuthGuard {
	return &AuthGuard{
		tokenService: tokenSvc,
	}
}

func (a *AuthGuard) Verify(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, ok := a.withAuthenticatedContext(r)

		if !ok {
			WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		h.ServeHTTP(w, req)
	})
}

func (a *AuthGuard) withAuthenticatedContext(r *http.Request) (*http.Request, bool) {
	ctx := r.Context()

	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, false
	}

	tokenStr := extractBearerToken(authorizationHeader)
	if tokenStr == "" {
		return nil, false
	}

	claim, err := a.tokenService.Validate(tokenStr)
	if err != nil {
		return nil, false
	}

	ctx = context.WithValue(ctx, ContextKeyUserID, claim.UserID)
	ctx = context.WithValue(ctx, ContextKeyEmail, claim.Email)

	return r.WithContext(ctx), true
}

func extractBearerToken(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.ToUpper(parts[0]) == "BEARER" {
		return parts[1]
	}
	return ""
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(string)
	return id, ok && id != ""
}
