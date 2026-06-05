// inbound (primary adapter)
package ports

import (
	"context"
	"io"
	"time"
)

type UploadImageCommand struct {
	OwnerID     string
	Filename    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type ImageUseCase interface {
	Upload(ctx context.Context, cmd UploadImageCommand) (ImageDTO, error)
}

type ImageDTO struct {
	ID            string    `json:"id"`
	OwnerID       string    `json:"owner_id"`
	OriginalName  string    `json:"original_name"`
	Format        string    `json:"format"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	FileSizeBytes int64     `json:"file_size_bytes"`
	Status        string    `json:"status"`
	IsPublic      bool      `json:"is_public"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	URL           string    `json:"url"`
}
