// inbound (primary adapter)

package ports

import (
	"context"
)

type RegisterCommand struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string
}

type UserUseCase interface {
	Register(ctx context.Context, cmd RegisterCommand) (AuthResult, error)
}
