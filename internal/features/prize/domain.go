package prize

import (
	"context"
	"fmt"
	"spinLuck/internal/shared/models"
	"strings"
)

type PrizeRepository interface {
	WithTransaction(fn func(repo PrizeRepository) error) error
	AssignToRaffle(ctx context.Context, raffleID uint64, prizeID uint64) error
	GetAllByIdRaffleByOrganizerID(ctx context.Context, raffleID uint64, organizerID uint64) ([]models.Prize, error)
	GetAllByIdRafflePublic(ctx context.Context, slug string) ([]models.Prize, error)

	GetByID(ctx context.Context, id uint64) (*models.Prize, error)
	ExistsByIDAndOrganizerID(ctx context.Context, id uint64, organizerID uint64) (bool, error)

	Create(ctx context.Context, prize *models.Prize) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	Delete(ctx context.Context, id uint64) error
	DeleteFromRaffle(ctx context.Context, raffleID uint64, prizeID uint64) error
}

type PrizeService interface {
	GetAllByIdRaffleByOrganizerID(ctx context.Context, raffleID uint64, organizerID uint64) ([]models.Prize, error)
	GetAllByIdRafflePublic(ctx context.Context, slug string) ([]models.Prize, error)
	GetByID(ctx context.Context, id uint64) (*models.Prize, error)

	Create(ctx context.Context, userID uint64, title string, description string, imageURL string, raffleID uint64) (*models.Prize, error)
	Update(ctx context.Context, userID uint64, id uint64, title, description, imageURL string) (*models.Prize, error)
	Delete(ctx context.Context, userID uint64, raffleID uint64, prizeID uint64) error
}

func NewPrize(title, description, imageURL string, organizerID uint64) (*models.Prize, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("el título del premio no puede estar vacío")
	}

	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("la descripción del premio no puede estar vacía")
	}

	if strings.TrimSpace(imageURL) == "" {
		return nil, fmt.Errorf("la URL de la imagen del premio no puede estar vacía")
	}

	if organizerID == 0 {
		return nil, fmt.Errorf("el ID del organizador es inválido")
	}

	return &models.Prize{
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		OrganizerID: organizerID,
	}, nil
}

func BuildPrizeUpdateData(title, description, imageURL string) (map[string]any, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("el título del premio no puede estar vacío")
	}

	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("la descripción del premio no puede estar vacía")
	}

	if strings.TrimSpace(imageURL) == "" {
		return nil, fmt.Errorf("la URL de la imagen del premio no puede estar vacía")
	}

	return map[string]any{
		"title":       title,
		"description": description,
		"image_url":   imageURL,
	}, nil
}
