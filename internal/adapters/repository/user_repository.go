package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (repo *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`

	u, err := scanUser(repo.db.QueryRowContext(ctx, query, email.String()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", email, shared.ErrNotFound)
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return u, nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`

	u, err := scanUser(repo.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", id, shared.ErrNotFound)
		}
		return nil, fmt.Errorf("find user by ID: %w", err)
	}

	return u, nil
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

func scanUser(row *sql.Row) (*user.User, error) {
	var (
		id           string
		emailStr     string
		passwordHash []byte
		createdAt    time.Time
		updatedAt    time.Time
	)

	if err := row.Scan(
		&id, &emailStr, &passwordHash, &createdAt, &updatedAt,
	); err != nil {
		return nil, err
	}

	emailVO, err := user.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("reconstitute user: bad email in db %q: %w", emailStr, err)
	}

	return user.ReconstitueUser(
		user.UserID(id),
		emailVO,
		passwordHash,
		createdAt,
		updatedAt,
	), nil
}
