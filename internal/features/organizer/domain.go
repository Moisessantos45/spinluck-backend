package organizer

import (
	"context"
	"fmt"
	"spinLuck/internal/shared/models"
	"strings"
)

type OrganizerBasicInfo struct {
	ID     uint64 `gorm:"column:id" json:"id"`
	UserID uint64 `gorm:"column:user_id" json:"user_id"`
}

type OrganizerDashboardMetrics struct {
	TotalAmount                float64 `json:"total_amount" gorm:"column:total_amount"`
	AverageTicketPrice         float64 `json:"average_ticket_price" gorm:"column:average_ticket_price"`
	TotalParticipants          uint64  `json:"total_participants" gorm:"column:total_participants"`
	TotalRaffles               uint64  `json:"total_raffles" gorm:"column:total_raffles"`
	TotalUniqueParticipants    uint64  `json:"total_unique_participants" gorm:"column:total_unique_participants"`
	EffectiveSalesRate         float64 `json:"effective_sales_rate" gorm:"column:effective_sales_rate"`
	ActiveRafflesExpiringSoon  uint64  `json:"active_raffles_expiring_soon" gorm:"column:active_raffles_expiring_soon"`
	EstimatedDepletionDays     float64 `json:"estimated_depletion_days" gorm:"column:estimated_depletion_days"`
	StagnantTickets            uint64  `json:"stagnant_tickets" gorm:"column:stagnant_tickets"`
	ActiveRafflesWithoutPrizes uint64  `json:"active_raffles_without_prizes" gorm:"column:active_raffles_without_prizes"`
}

type OrganizerRepository interface {
	GetDashboardInfoDetailGeneric(ctx context.Context, organizerID uint64) (OrganizerDashboardMetrics, error)
	GetByUserID(userID uint64) (*models.Organizer, error)
	GetByID(id uint64) (*models.Organizer, error)
	GetBasicInfoByUserID(userID uint64) (*OrganizerBasicInfo, error)
	Create(organizer *models.Organizer) error
	Update(id uint64, data map[string]any) error
}

type OrganizerService interface {
	GetDashboardInfoDetailGeneric(ctx context.Context, userID uint64) (OrganizerDashboardMetrics, error)
	GetByUserID(userID uint64) (*models.Organizer, error)
	GetBasicInfoByUserID(userID uint64) (*OrganizerBasicInfo, error)
	GetByID(id uint64) (*models.Organizer, error)
	Create(name string, phone string, userID uint64) (*models.Organizer, error)
	Update(id uint64, name string, phone string) (*models.Organizer, error)
}

func NewOrganizer(name string, phone string, userID uint64) (*models.Organizer, error) {
	if userID == 0 {
		return nil, fmt.Errorf("El ID del usuario no puede ser cero")
	}

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("El nombre del organizador no puede estar vacío")
	}
	if strings.TrimSpace(phone) == "" {
		return nil, fmt.Errorf("El teléfono del organizador no puede estar vacío")
	}

	phoneNotSpace := strings.ReplaceAll(phone, " ", "")

	return &models.Organizer{
		Name:   name,
		Phone:  phoneNotSpace,
		UserID: userID,
	}, nil
}

func BuildOrganizerUpdateData(name string, phone string) (map[string]any, error) {

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("El nombre del organizador no puede estar vacío")
	}

	if strings.TrimSpace(phone) == "" {
		return nil, fmt.Errorf("El teléfono del organizador no puede estar vacío")
	}

	phoneNotSpace := strings.ReplaceAll(phone, " ", "")

	return map[string]any{
		"name":  name,
		"phone": phoneNotSpace,
	}, nil
}
