package ports

import "context"

type StoragePath struct {
	Path string
	URL  string
}

type ImageStorage interface {
	Store(ctx context.Context, data []byte, key, mimeType string) (StoragePath, error)
}
