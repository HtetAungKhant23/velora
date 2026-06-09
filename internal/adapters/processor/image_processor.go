package processor

import (
	"bytes"
	"context"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/webp"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// ctx will need later for cancellation/timeout support
func (p *ImageProcessor) DecodeMetadata(_ context.Context, data []byte) (int, int, error) {
	img, err := decodeImage(data)
	if err != nil {
		return 0, 0, fmt.Errorf("processor: decode metadata: %w", err)
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

func decodeImage(data []byte) (image.Image, error) {
	r := bytes.NewReader(data)

	if img, err := webp.Decode(r); err == nil {
		return img, nil
	}
	r.Seek(0, 0) // reset reader position

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return img, nil
}
