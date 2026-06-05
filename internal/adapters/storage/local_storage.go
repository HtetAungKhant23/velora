package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HtetAungKhant23/velora/internal/core/ports"
)

const (
	mkdirPerm     = 0o755
	writeFilePerm = 0o644
)

type LocalStorage struct {
	baseDir string
	baseURL string
}

func NewLocalStorage(baseDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, mkdirPerm); err != nil {
		return nil, fmt.Errorf("storage: create base dir: %w", err)
	}

	return &LocalStorage{
		baseDir: baseDir,
		baseURL: baseURL,
	}, nil
}

// other implementations will need ctx and mimetype
func (s *LocalStorage) Store(_ context.Context, data []byte, key string, _ string) (ports.StoragePath, error) {
	fullpath := filepath.Join(s.baseDir, filepath.FromSlash(key))

	if err := os.MkdirAll(filepath.Dir(fullpath), mkdirPerm); err != nil {
		return ports.StoragePath{}, fmt.Errorf("storage: mkdir: %w", err)
	}

	if err := os.WriteFile(fullpath, data, writeFilePerm); err != nil {
		return ports.StoragePath{}, fmt.Errorf("storage: write file: %w", err)
	}

	return ports.StoragePath{
		Path: key,
		URL:  s.PublicURL(key),
	}, nil
}

func (s *LocalStorage) PublicURL(path string) string {
	return s.baseURL + "/" + path
}
