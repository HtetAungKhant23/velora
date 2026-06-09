package handler

import (
	"errors"
	"net/http"

	"github.com/HtetAungKhant23/velora/internal/adapters/handler/middleware"
	imageDom "github.com/HtetAungKhant23/velora/internal/core/domain/image"
	"github.com/HtetAungKhant23/velora/internal/core/ports"
)

const maxUploadSize = 50 * 1024 * 1024

type ImageHandler struct {
	usecase ports.ImageUseCase
}

func NewImageHandler(uc ports.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		usecase: uc,
	}
}

func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "file too large or invliad multipart form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "field 'image' required")
		return
	}
	defer file.Close()

	uploadCmd := ports.UploadImageCommand{
		OwnerID:     userID,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		Reader:      file,
	}

	imageDTO, err := h.usecase.Upload(ctx, uploadCmd)
	if err != nil {
		switch {
		case errors.Is(err, imageDom.ErrUnsupportedFormat):
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, imageDom.ErrFileSizeTooLarge):
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, imageDom.ErrFileSizeZero):
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, imageDom.ErrImageStillUploading):
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
			return
		default:
			middleware.MapDomainError(w, err)
			return
		}
	}

	middleware.WriteCreated(w, imageDTO)
}
