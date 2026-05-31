package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/HtetAungKhant23/velora/internal/core/domain/shared"
	"github.com/HtetAungKhant23/velora/internal/core/domain/user"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	if err := repo.db.QueryRowContext(ctx, query, email.String()).Scan(&exists); err != nil {
		return false, fmt.Errorf("exists by email: %w", err)
	}

	return exists, nil
}

func (repo *UserRepository) Save(ctx context.Context, user *user.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`

	_, err := repo.db.ExecContext(ctx, query, user.Email().String(), user.PasswordHash())
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("save user: %w", shared.ErrAlreadyExist)
		}
		return fmt.Errorf("save user: %w", err)
	}

	return nil
}
