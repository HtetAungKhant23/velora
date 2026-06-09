package image

import (
	"errors"
	"strconv"
	"time"

	"github.com/HtetAungKhant23/velora/internal/core/domain/user"
	"github.com/google/uuid"
)

type ImageID string

func newImageID() ImageID {
	return ImageID(uuid.New().String())
}

type ImageStatus string

const (
	StatusUploading  ImageStatus = "uploading"
	StatusReady      ImageStatus = "ready"
	StatusProcessing ImageStatus = "processing"
	StatusDeleted    ImageStatus = "deleted"
)

var (
	ErrImageStillUploading = errors.New("image: upload still in progress")
)

type Image struct {
	id           ImageID
	ownerID      user.UserID
	originalName string
	format       Format
	dimensions   Dimensions
	status       ImageStatus
	storagePath  StoragePath
	isPublic     bool
	fileSize     FileSize
	createdAt    time.Time
	updatedAt    time.Time
}

func NewImage(
	ownerID user.UserID,
	originalName string,
	format Format,
	fileSize FileSize,
	storagePath StoragePath,
) (*Image, error) {
	if ownerID == "" {
		return nil, errors.New("image: ownerID is required")
	}

	if originalName == "" {
		return nil, errors.New("image: original name is required")
	}

	now := time.Now().UTC()
	img := &Image{
		id:           newImageID(),
		ownerID:      ownerID,
		originalName: originalName,
		format:       format,
		fileSize:     fileSize,
		storagePath:  storagePath,
		status:       StatusUploading,
		createdAt:    now,
		updatedAt:    now,
	}

	return img, nil
}

func (img *Image) ID() ImageID {
	return img.id
}

func (img *Image) OwnerID() user.UserID {
	return img.ownerID
}

func (img *Image) OriginalName() string {
	return img.originalName
}

func (img *Image) Format() Format {
	return img.format
}

func (img *Image) Dimensions() Dimensions {
	return img.dimensions
}

func (img *Image) Status() ImageStatus {
	return img.status
}

func (img *Image) StoragePath() StoragePath {
	return img.storagePath
}

func (img *Image) IsPublic() bool {
	return img.isPublic
}

func (img *Image) FileSize() FileSize {
	return img.fileSize
}

func (img *Image) IsReady() bool {
	return img.status == StatusReady
}

func (img *Image) CreatedAt() time.Time {
	return img.createdAt
}

func (img *Image) UpdatedAt() time.Time {
	return img.updatedAt
}

func (img *Image) AssignID(id int64) {
	img.id = ImageID(strconv.FormatInt(id, 10))
}

func (img *Image) MarkReady(dims Dimensions) error {
	if img.status != StatusUploading {
		return ErrImageStillUploading
	}

	img.dimensions = dims
	img.status = StatusReady
	img.updatedAt = time.Now().UTC()

	return nil
}
