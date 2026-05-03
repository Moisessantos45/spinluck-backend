package ticket

import (
	"context"
	"fmt"
	"log"
	"spinLuck/internal/shared/models"
	"strings"
)

type TicketWithOrganizerNumber struct {
	ID              uint64 `gorm:"column:id" json:"id"`
	Number          uint64 `gorm:"column:number" json:"number"`
	Slug            string `gorm:"column:slug" json:"slug"`
	OrganizerNumber string `gorm:"column:organizer_number" json:"organizer_number"`
	FormattedNumber string `json:"formatted_number" gorm:"column:formatted_number"`
}

type TicketVoucherData struct {
	Slug            string `json:"slug" gorm:"column:slug"`
	TicketID        uint64 `json:"ticket_id" gorm:"column:ticket_id"`
	Number          uint64 `json:"number" gorm:"column:number"`
	FormattedNumber string `json:"formatted_number" gorm:"column:formatted_number"`
	Amount          uint64 `json:"amount" gorm:"column:amount"`
	FullName        string `json:"full_name" gorm:"column:full_name"`
}

type TicketRepository interface {
	GetAllStatus(ctx context.Context) ([]models.TicketStatus, error)
	GetAll(ctx context.Context) ([]models.Ticket, error)
	GetAllByRaffleIDToOrgaizer(ctx context.Context, raffleID uint64, organizerID uint64) ([]models.Ticket, error)

	IsTicketFromOrganizer(ctx context.Context, ticketID uint64, organizerID uint64) (bool, error)
	GetByID(ctx context.Context, id uint64, organizerID uint64) (*models.Ticket, error)
	GetVoucherDataByID(ctx context.Context, id uint64) (*TicketVoucherData, error)
	GetRandomSoldTicket(ctx context.Context, raffleID uint64) (*models.Ticket, error)
	GetRandomAvailableTicket(ctx context.Context, raffleID uint64) (*TicketWithOrganizerNumber, error)
	FindByNumberAndRaffleID(ctx context.Context, number uint64, raffleID uint64) (*TicketWithOrganizerNumber, error)

	Create(ctx context.Context, ticket *models.Ticket) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	UpdateWinner(ctx context.Context, id uint64, raffleID uint64) error
	Delete(ctx context.Context, id uint64) error
}

type TicketService interface {
	GetAllStatus(ctx context.Context) ([]models.TicketStatus, error)
	GetAll(ctx context.Context) ([]models.Ticket, error)
	GetAllByRaffleIDToOrgaizer(ctx context.Context, raffleID uint64, userID uint64) ([]models.Ticket, error)

	GetByID(ctx context.Context, id uint64, userID uint64) (*models.Ticket, error)
	GetRandomSoldTicket(ctx context.Context, userID uint64, raffleID uint64) (*models.Ticket, error)
	GetRandomAvailableTicket(ctx context.Context, raffleID uint64) (*TicketWithOrganizerNumber, error)
	FindByNumberAndRaffleID(ctx context.Context, number uint64, raffleID uint64) (*TicketWithOrganizerNumber, error)
	GenerateVoucher(ctx context.Context, ticketID uint64, userID uint64) ([]byte, error)

	Create(ctx context.Context, number uint64, participantName, participantPhone string, raffleID uint64) (*models.Ticket, error)
	Update(ctx context.Context, id uint64, participantName, participantPhone string, ticketStatusID uint64, userID uint64) (*models.Ticket, error)
	UpdateWinner(ctx context.Context, id uint64, raffleID uint64, userID uint64) error
	Delete(ctx context.Context, id uint64, userID uint64) error
}

func NewTicket(number uint64, participantName string, participantPhone string, raffleID uint64) (*models.Ticket, error) {
	if raffleID == 0 {
		return nil, fmt.Errorf("el ID de la rifa es inválido")
	}

	if strings.TrimSpace(participantName) == "" {
		return nil, fmt.Errorf("el nombre del participante es inválido")
	}

	if strings.TrimSpace(participantPhone) == "" {
		return nil, fmt.Errorf("el teléfono del participante es inválido")
	}

	phoneNotDigits := strings.ReplaceAll(participantPhone, " ", "")

	return &models.Ticket{
		Number:           number,
		ParticipantName:  participantName,
		ParticipantPhone: phoneNotDigits,
		RaffleID:         raffleID,
		TicketStatusID:   2,
	}, nil
}

func BuildTicketUpdateData(participantName, participantPhone string, ticketStatusID uint64) (map[string]any, error) {

	if ticketStatusID == 0 {
		return nil, fmt.Errorf("el ID del estado del ticket es inválido")
	}

	if strings.TrimSpace(participantName) == "" {
		return nil, fmt.Errorf("el nombre del participante es inválido")
	}

	if strings.TrimSpace(participantPhone) == "" {
		return nil, fmt.Errorf("el teléfono del participante es inválido")
	}

	phoneNotDigits := strings.ReplaceAll(participantPhone, " ", "")
	log.Printf("Ticket update data - participantName: %s, participantPhone: %s, ticketStatusID: %d", participantName, phoneNotDigits, ticketStatusID)

	return map[string]any{
		"participant_name":  participantName,
		"participant_phone": phoneNotDigits,
		"ticket_status_id":  ticketStatusID,
	}, nil
}
