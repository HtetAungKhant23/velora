package ports

import (
	"context"

	imageDom "github.com/HtetAungKhant23/velora/internal/core/domain/image"
)

type ImageRepository interface {
	Save(ctx context.Context, img *imageDom.Image) error
}
