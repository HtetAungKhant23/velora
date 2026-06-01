// do NOT contain business logic — that lives in the User entity.

package services

import (
	"context"
	"errors"
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

func (s *UserService) Register(ctx context.Context, cmd ports.RegisterCommand) error {
	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	exist, err := s.users.ExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("register: check email: %w", err)
	}

	if exist {
		return fmt.Errorf("register: %w", shared.ErrAlreadyExist)
	}

	user, err := user.NewUser(cmd.Email, cmd.Password)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	if err = s.users.Save(ctx, user); err != nil {
		return fmt.Errorf("register: save: %w", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, cmd ports.LoginCommand) (ports.AuthResult, error) {
	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return ports.AuthResult{}, shared.ErrUnauthorized
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return ports.AuthResult{}, shared.ErrUnauthorized
		}
		return ports.AuthResult{}, fmt.Errorf("login: %w", err)
	}

	if err = user.VerifyPassword(cmd.Password); err != nil {
		return ports.AuthResult{}, shared.ErrUnauthorized
	}

	token, err := s.tokens.Generate(string(user.ID()), user.Email().String())
	if err != nil {
		return ports.AuthResult{}, fmt.Errorf("login: token: %w", err)
	}

	return ports.AuthResult{
		Token: token,
	}, nil
}
