// do NOT contain business logic — that lives in the User entity.

package services

import (
	"context"
	"fmt"

	"github.com/HtetAungKhant23/velora/internal/core/domain/shared"
	"github.com/HtetAungKhant23/velora/internal/core/domain/user"
	"github.com/HtetAungKhant23/velora/internal/core/ports"
)

type UserService struct {
	users  ports.UserRepository
	tokens ports.TokenService
}

func NewUserService(users ports.UserRepository, tokens ports.TokenService) *UserService {
	return &UserService{
		users:  users,
		tokens: tokens,
	}
}

func (s *UserService) Register(ctx context.Context, cmd ports.RegisterCommand) (ports.AuthResult, error) {
	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return ports.AuthResult{}, fmt.Errorf("register: %w", err)
	}

	exist, err := s.users.ExistsByEmail(ctx, email)
	if err != nil {
		return ports.AuthResult{}, fmt.Errorf("register: check email: %w", err)
	}

	if exist {
		return ports.AuthResult{}, fmt.Errorf("register: %w", shared.ErrAlreadyExist)
	}

	user, err := user.NewUser(cmd.Email, cmd.Password)
	if err != nil {
		return ports.AuthResult{}, fmt.Errorf("register: %w", err)
	}

	if err = s.users.Save(ctx, user); err != nil {
		return ports.AuthResult{}, fmt.Errorf("register: save: %w", err)
	}

	token, err := s.tokens.Generate(string(user.ID()), user.Email().String())
	if err != nil {
		return ports.AuthResult{}, fmt.Errorf("register: generate token: %w", err)
	}

	return ports.AuthResult{
		Token: token,
	}, nil
}
