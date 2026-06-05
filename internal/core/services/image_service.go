package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	imageDom "github.com/HtetAungKhant23/velora/internal/core/domain/image"
	"github.com/HtetAungKhant23/velora/internal/core/domain/user"
	"github.com/HtetAungKhant23/velora/internal/core/ports"
)

type ImageService struct{}

func NewImageService() *ImageService {
	return &ImageService{}
}

func (s *ImageService) Upload(ctx context.Context, cmd ports.UploadImageCommand) (ports.ImageDTO, error) {
	format, err := s.resolveFormat(cmd.Filename, cmd.ContentType)
	if err != nil {
		return ports.ImageDTO{}, err
	}

	fileSize, err := imageDom.NewFileSize(cmd.Size)
	if err != nil {
		return ports.ImageDTO{}, err
	}

	storageKey := buildStorageKey(cmd.OwnerID, cmd.Filename)
	storagePath, err := imageDom.NewStoragePath(storageKey)
	if err != nil {
		return ports.ImageDTO{}, fmt.Errorf("upload: storage path: %w", err)
	}

	img, err := imageDom.NewImage(
		user.UserID(cmd.OwnerID),
		cmd.Filename,
		format,
		fileSize,
		storagePath,
	)
	if err != nil {
		return ports.ImageDTO{}, fmt.Errorf("upload: image: %w", err)
	}

	return toImageDTO(img), nil
}

func (s *ImageService) resolveFormat(filename, contentType string) (imageDom.Format, error) {
	f, err := imageDom.FormatFromFilename(filename)
	if err == nil {
		return f, nil
	}

	// fallback from content-type
	ct := strings.ToLower(contentType)
	switch {
	case strings.Contains(ct, "jpeg"):
		return imageDom.FormatJPEG, nil
	case strings.Contains(ct, "jpg"):
		return imageDom.FormatJPG, nil
	case strings.Contains(ct, "png"):
		return imageDom.FormatPNG, nil
	case strings.Contains(ct, "gif"):
		return imageDom.FormatGIF, nil
	case strings.Contains(ct, "webp"):
		return imageDom.FormatWEBP, nil
	}

	return imageDom.Format{}, err
}

func toImageDTO(img *imageDom.Image) ports.ImageDTO {
	return ports.ImageDTO{
		ID:            string(img.ID()),
		OwnerID:       string(img.OwnerID()),
		OriginalName:  img.OriginalName(),
		Format:        img.Format().String(),
		Width:         img.Dimensions().Width(),
		Height:        img.Dimensions().Height(),
		FileSizeBytes: img.FileSize().Bytes(),
		Status:        string(img.Status()),
		IsPublic:      img.IsPublic(),
		CreatedAt:     img.CreatedAt(),
		UpdatedAt:     img.UpdatedAt(),
	}
}

func buildStorageKey(ownerID, filename string) string {
	return fmt.Sprintf("images/%s/%s/%s",
		ownerID,
		strings.ReplaceAll(filename, " ", "_"),
		fmt.Sprintf("%d", time.Now().UnixNano()),
	)
}
