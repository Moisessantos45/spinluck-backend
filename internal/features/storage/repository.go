package storage

import (
	"context"
	"errors"
	"spinLuck/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) StorageRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) WithTransaction(fn func(repo StorageRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &PostgresRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *PostgresRepository) GetByFileId(ctx context.Context, fileId string) (*models.File, error) {
	var vacante models.File
	err := r.db.Where("file_id = ?", fileId).First(&vacante).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	return &vacante, nil
}

func (r *PostgresRepository) Create(ctx context.Context, file *models.File) error {
	return r.db.Create(file).Error
}

func (r *PostgresRepository) Update(ctx context.Context, userID uint64, file *models.File) error {

	updates := map[string]any{
		"path":      file.Path,
		"mime_type": file.MimeType,
		"size":      file.Size,
	}

	return r.db.Model(&models.File{}).Where("file_id = ? AND user_id = ?", file.FileID, userID).Updates(updates).Error
}

func (r *PostgresRepository) Delete(ctx context.Context, userID uint64, fileId string) error {
	return r.db.Where("file_id = ? AND user_id = ?", fileId, userID).Delete(&models.File{}).Error
}
