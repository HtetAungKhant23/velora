package image

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	FormatJPG  = Format{"jpg"}
	FormatJPEG = Format{"jpeg"}
	FormatPNG  = Format{"png"}
	FormatGIF  = Format{"gif"}
	FormatWEBP = Format{"webp"}

	ErrUnsupportedFormat = errors.New("format: unsupported(jpg|jpeg|png|webp|gif)")
)

var validFormat = map[string]Format{
	"jpg":  FormatJPG,
	"jpeg": FormatJPEG,
	"png":  FormatPNG,
	"gif":  FormatGIF,
	"webp": FormatWEBP,
}

type Format struct {
	value string
}

func NewFormat(raw string) (Format, error) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if format, ok := validFormat[raw]; ok {
		return format, nil
	}
	return Format{}, ErrUnsupportedFormat
}

func (f Format) String() string {
	return f.value
}

func (f Format) Equals(format Format) bool {
	return f.value == format.value
}

func FormatFromFilename(filename string) (Format, error) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	return NewFormat(ext)
}

func (f Format) MimeType() string {
	switch f.value {
	case "jpeg":
		return "image/jpeg"
	case "jpg":
		return "image/jpg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		return "application/octect-stream"
	}
}
