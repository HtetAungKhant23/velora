package ports

import "context"

type ImageProcessor interface {
	DecodeMetadata(ctx context.Context, data []byte) (int, int, error)
}
