package repository

import (
	"context"
	"database/sql"

	"github.com/HtetAungKhant23/velora/internal/core/domain/user"
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
	return false, nil
}

func (repo *UserRepository) Save(ctx context.Context, user *user.User) error {
	return nil
}
