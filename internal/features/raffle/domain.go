package raffle

import (
	"context"
	"fmt"
	"spinLuck/internal/shared/models"
	"strings"
	"time"
)

type RaffleBasicInfo struct {
	ID          uint64 `gorm:"column:id" json:"id"`
	OrganizerID uint64 `gorm:"column:organizer_id" json:"organizer_id"`
	ImageURL    string `gorm:"column:image_url" json:"image_url"`
	Quantity    uint64 `gorm:"column:quantity_tickets" json:"quantity_tickets"`
}

type RaffleInfoGeneric struct {
	ID              uint64  `json:"id" gorm:"column:id"`
	Title           string  `json:"title" gorm:"column:title"`
	Slug            string  `json:"slug" gorm:"column:slug"`
	QuantityTickets uint64  `json:"quantity_tickets" gorm:"column:quantity_tickets"`
	Status          string  `json:"status" gorm:"column:status"`
	ImageURL        string  `json:"image_url" gorm:"column:image_url"`
	Progress        float64 `json:"progress" gorm:"column:progress"`
	TotalSold       float64 `json:"total_sold" gorm:"column:total_sold"`
	TotalAmount     float64 `json:"total_amount" gorm:"column:total_amount"`
}

type RaffleRepository interface {
	WithTransaction(fn func(repo RaffleRepository) error) error

	GetAllStatus(ctx context.Context) ([]models.RaffleStatus, error)
	GetAll(ctx context.Context) ([]models.Raffle, error)
	GetAllInfoGeneric(ctx context.Context, organizerID uint64) ([]RaffleInfoGeneric, error)
	GetAllRecentInfoGeneric(ctx context.Context, organizerID uint64) ([]RaffleInfoGeneric, error)

	GetByID(ctx context.Context, id uint64) (*models.Raffle, error)
	GetByOrganizerID(ctx context.Context, organizerID uint64) ([]models.Raffle, error)

	GetInfoBasicByID(ctx context.Context, id uint64, organizerID uint64) (*RaffleBasicInfo, error)
	GetBySlug(ctx context.Context, slug string) (*models.Raffle, error)

	Create(ctx context.Context, raffle *models.Raffle) error
	CreateTickets(ctx context.Context, tickets []models.Ticket) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	GetByIDForUpdate(ctx context.Context, id uint64) (*models.Raffle, error)

	CountNonAvailableTicketsFromNumber(ctx context.Context, raffleID uint64, fromNumber uint64) (int64, error)
	DeleteAvailableTicketsFromNumber(ctx context.Context, raffleID uint64, fromNumber uint64) error
}

type RaffleService interface {
	GetAllStatus(ctx context.Context) ([]models.RaffleStatus, error)
	GetAll(ctx context.Context) ([]models.Raffle, error)
	GetAllInfoGeneric(ctx context.Context, userID uint64) ([]RaffleInfoGeneric, error)
	GetAllRecentInfoGeneric(ctx context.Context, userID uint64) ([]RaffleInfoGeneric, error)

	GetByID(ctx context.Context, id uint64) (*models.Raffle, error)
	GetBySlug(ctx context.Context, userID uint64, slug string) (*models.Raffle, error)
	GetByOrganizerID(ctx context.Context, organizerID uint64) ([]models.Raffle, error)

	Create(ctx context.Context, userID uint64, title string, description string, price float64, imageURL string, dateStr string, quantityTickets uint64, maxWinners uint64, raffleStatusID uint64) (*models.Raffle, error)
	Update(ctx context.Context, userID uint64, id uint64, title, description string, price float64, imageURL, dateStr string, quantityTickets uint64, maxWinners uint64, raffleStatusID uint64) (*models.Raffle, error)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	return slug
}

func NewRaffle(title, description string, price float64, imageURL, date string, quantityTickets uint64, maxWinners uint64, organizerID uint64, raffleStatusID uint64) (*models.Raffle, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("el título no puede estar vacío")
	}

	if price <= 0 {
		return nil, fmt.Errorf("el precio debe ser mayor a 0")
	}

	if quantityTickets <= 0 {
		return nil, fmt.Errorf("la cantidad de tickets debe ser mayor a 0")
	}

	if maxWinners <= 0 {
		return nil, fmt.Errorf("la cantidad máxima de ganadores debe ser mayor a 0")
	}

	if organizerID == 0 {
		return nil, fmt.Errorf("el ID del organizador es inválido")
	}

	if strings.TrimSpace(date) == "" {
		return nil, fmt.Errorf("la fecha no puede estar vacía")
	}

	if raffleStatusID == 0 {
		return nil, fmt.Errorf("el ID del estado de la rifa es inválido")
	}

	return &models.Raffle{
		Title:           title,
		Description:     description,
		Price:           price,
		ImageURL:        imageURL,
		Date:            date,
		QuantityTickets: quantityTickets,
		MaxWinners:      maxWinners,
		Slug:            generateSlug(title),
		OrganizerID:     organizerID,
		RaffleStatusID:  raffleStatusID,
	}, nil
}

func BuildRaffleUpdateData(
	title string,
	description string,
	price float64,
	imageURL string,
	date string,
	quantityTickets uint64,
	maxWinners uint64,
	raffleStatusID uint64,
) (map[string]any, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("el título no puede estar vacío")
	}

	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("la descripción no puede estar vacía")
	}

	if price <= 0 {
		return nil, fmt.Errorf("el precio debe ser mayor a 0")
	}

	if quantityTickets == 0 {
		return nil, fmt.Errorf("la cantidad de tickets debe ser mayor a 0")
	}

	if maxWinners == 0 {
		return nil, fmt.Errorf("la cantidad máxima de ganadores debe ser mayor a 0")
	}

	if strings.TrimSpace(date) == "" {
		return nil, fmt.Errorf("la fecha no puede estar vacía")
	}

	if raffleStatusID == 0 {
		return nil, fmt.Errorf("el ID del estado de la rifa es inválido")
	}

	return map[string]any{
		"title":            strings.TrimSpace(title),
		"description":      strings.TrimSpace(description),
		"price":            price,
		"image_url":        strings.TrimSpace(imageURL),
		"date":             date,
		"quantity_tickets": quantityTickets,
		"max_winners":      maxWinners,
		"raffle_status_id": raffleStatusID,
	}, nil
}
