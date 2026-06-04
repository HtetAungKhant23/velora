package image

import (
	"errors"
	"strings"
)

var ErrStoragePathEmpty = errors.New("storage path: cannot be empty")

type StoragePath struct {
	value string
}

func NewStoragePath(path string) (StoragePath, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return StoragePath{}, ErrStoragePathEmpty
	}
	return StoragePath{path}, nil
}

func (s StoragePath) String() string {
	return s.value
}

func (s StoragePath) Equals(o StoragePath) bool {
	return s.value == o.value
}
