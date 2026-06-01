// inbound (primary adapter)

package ports

import (
	"context"
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

type UserUseCase interface {
	Register(ctx context.Context, cmd RegisterCommand) error
	Login(ctx context.Context, cmd LoginCommand) (AuthResult, error)
}
