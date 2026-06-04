package image

import "errors"

var (
	ErrDimensionZero     = errors.New("dimension: width and height must be > 0")
	ErrDimensionTooLarge = errors.New("dimension: maximum 16000 x 16000 pixels")
)

const maxDimension = 16000

type Dimensions struct {
	width  int
	height int
}

func NewDimensions(width, height int) (Dimensions, error) {
	if width <= 0 || height <= 0 {
		return Dimensions{}, ErrDimensionZero
	}
	if width > maxDimension || height > maxDimension {
		return Dimensions{}, ErrDimensionTooLarge
	}
	return Dimensions{
		width:  width,
		height: height,
	}, nil
}

func (d Dimensions) Width() int {
	return d.width
}

func (d Dimensions) Height() int {
	return d.height
}

func (d Dimensions) Equals(o Dimensions) bool {
	return d.height == o.height && d.width == o.width
}
