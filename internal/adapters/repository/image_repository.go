package repository

import (
	"context"
	"database/sql"
	"fmt"

	imageDom "github.com/HtetAungKhant23/velora/internal/core/domain/image"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

func (repo *ImageRepository) Save(ctx context.Context, img *imageDom.Image) error {
	query := `
		INSERT INTO images
  			(owner_id, original_name, format,
     		width, height, file_size, storage_path,
       		status, is_public, created_at, updated_at)
        VALUES (
            $1, $2, $3,
            $4, $5, $6, $7,
            $8, $9, $10, $11
        )
        RETURNING id
	`

	var id int64
	err := repo.db.QueryRowContext(ctx, query,
		img.OwnerID(),
		img.OriginalName(),
		img.Format().String(),
		img.Dimensions().Width(),
		img.Dimensions().Height(),
		img.FileSize().Bytes(),
		img.StoragePath().String(),
		string(img.Status()),
		img.IsPublic(),
		img.CreatedAt(),
		img.UpdatedAt(),
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	img.AssignID(id)

	return nil
}
