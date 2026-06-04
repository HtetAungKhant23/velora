package image

import "errors"

var (
	ErrFileSizeZero     = errors.New("file size: must be > 0")
	ErrFileSizeTooLarge = errors.New("file size: exceeds 50 MB limit")
)

const maxFileSizeByte = 50 * 1024 * 1024

type FileSize struct {
	bytes int64
}

func NewFileSize(bytes int64) (FileSize, error) {
	if bytes <= 0 {
		return FileSize{}, ErrFileSizeZero
	}
	if bytes > maxFileSizeByte {
		return FileSize{}, ErrFileSizeTooLarge
	}

	return FileSize{bytes}, nil
}

func (s FileSize) Bytes() int64 {
	return s.bytes
}

func (s FileSize) MB() float64 {
	return float64(s.bytes) / (1024 * 1024)
}
