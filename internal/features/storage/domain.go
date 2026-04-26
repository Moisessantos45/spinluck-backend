package storage

import (
	"context"
	"fmt"
	"spinLuck/internal/shared/models"
	"strings"
)

type StorageRepository interface {
	WithTransaction(fn func(repo StorageRepository) error) error
	GetByFileId(ctx context.Context, fileId string) (*models.File, error)
	Create(ctx context.Context, file *models.File) error
	Update(ctx context.Context, userID uint64, file *models.File) error
	Delete(ctx context.Context, userID uint64, fileId string) error
}

type StorageServiceInterface interface {
	GetByFileId(ctx context.Context, fileId string) (*models.File, error)
	Create(ctx context.Context, userID uint64, fileId string, path string, mimeType string, size int64) (*models.File, error)
	Update(ctx context.Context, userID uint64, fileId string, path string, mimeType string, size int64) (*models.File, error)
	Delete(ctx context.Context, userID uint64, fileId string) error
}

func NewStorageService(userID uint64, fileId string, path string, mimeType string, size int64) (*models.File, error) {
	if userID == 0 {
		return nil, fmt.Errorf("El userID es obligatorio")
	}

	if fileId == "" || strings.TrimSpace(fileId) == "" {
		return nil, fmt.Errorf("El fileId es obligatorio")
	}

	if path == "" || strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("El path es obligatorio")
	}

	if mimeType == "" || strings.TrimSpace(mimeType) == "" || !strings.Contains(mimeType, "/") || (!strings.Contains(mimeType, "image/") && !strings.Contains(mimeType, "application/pdf")) {
		return nil, fmt.Errorf("El mimeType es obligatorio y debe ser una imagen o un PDF")
	}

	if size <= 0 {
		return nil, fmt.Errorf("El size es obligatorio")
	}

	return &models.File{
		FileID:   fileId,
		Path:     path,
		MimeType: mimeType,
		Size:     size,
		UserID:   userID,
	}, nil
}
