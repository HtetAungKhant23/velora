// inbound (primary adapter)

package ports

import (
	"context"
	"time"
)

type RegisterCommand struct {
	Email    string
	Password string
}

type LoginCommand struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string `json:"token"`
}

type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUseCase interface {
	Register(ctx context.Context, cmd RegisterCommand) error
	Login(ctx context.Context, cmd LoginCommand) (AuthResult, error)
	GetProfile(ctx context.Context, userID string) (UserDTO, error)
}
