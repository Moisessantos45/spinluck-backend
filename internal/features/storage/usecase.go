package storage

import (
	"context"
	"fmt"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"
)

type StorageUseCase struct {
	repo StorageRepository
}

func NewStorageUseCase(repo StorageRepository) StorageServiceInterface {
	return &StorageUseCase{repo: repo}
}

func (s *StorageUseCase) GetByFileId(ctx context.Context, fileId string) (*models.File, error) {
	return s.repo.GetByFileId(ctx, fileId)
}

func (s *StorageUseCase) Create(ctx context.Context, userID uint64, fileId string, path string, mimeType string, size int64) (*models.File, error) {
	file, err := NewStorageService(userID, fileId, path, mimeType, size)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el File: %v", err)
	}

	url := utils.GenerateURL("storage/file", fileId)

	file.Url = url
	return file, nil

}

func (s *StorageUseCase) Update(ctx context.Context, userID uint64, fileId string, path string, mimeType string, size int64) (*models.File, error) {
	file, err := NewStorageService(userID, fileId, path, mimeType, size)
	if err != nil {
		return nil, err
	}

	err = s.repo.Update(ctx, userID, file)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar el File: %v", err)
	}

	url := utils.GenerateURL("storage/file", fileId)

	file.Url = url

	return file, nil
}

func (s *StorageUseCase) Delete(ctx context.Context, userID uint64, fileId string) error {
	file, err := s.repo.GetByFileId(ctx, fileId)
	if err != nil {
		return fmt.Errorf("Error al obtener el File: %v", err)
	}

	err = s.repo.WithTransaction(func(repo StorageRepository) error {
		err := utils.DeleteFile(file.Path)
		if err != nil {
			return fmt.Errorf("Error al eliminar el archivo: %v", err)
		}

		err = repo.Delete(ctx, userID, fileId)
		if err != nil {
			return fmt.Errorf("Error al eliminar el registro en la base de datos: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
