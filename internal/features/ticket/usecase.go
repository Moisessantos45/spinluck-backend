package ticket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"
	"time"
)

type TicketUseCase struct {
	rpOrg organizer.OrganizerRepository
	repo  TicketRepository
}

func NewTicketUseCase(repo TicketRepository, ucOrg organizer.OrganizerRepository) TicketService {
	return &TicketUseCase{repo: repo, rpOrg: ucOrg}
}

func (uc *TicketUseCase) GetAllStatus(ctx context.Context) ([]models.TicketStatus, error) {
	return uc.repo.GetAllStatus(ctx)
}

func (uc *TicketUseCase) GetAll(ctx context.Context) ([]models.Ticket, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *TicketUseCase) GetAllByRaffleIDToOrgaizer(ctx context.Context, raffleID uint64, userID uint64) ([]models.Ticket, error) {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	return uc.repo.GetAllByRaffleIDToOrgaizer(ctx, raffleID, organizer.ID)
}

func (uc *TicketUseCase) GetByID(ctx context.Context, id uint64, userID uint64) (*models.Ticket, error) {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	return uc.repo.GetByID(ctx, id, organizer.ID)
}

func (uc *TicketUseCase) GetRandomSoldTicket(ctx context.Context, userID uint64, raffleID uint64) (*models.Ticket, error) {
	organizer, err := uc.rpOrg.GetBasicInfoByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	exists, err := uc.repo.IsTicketFromOrganizer(ctx, raffleID, organizer.ID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("no se encontraron tickets para esta rifa y organizador")
	}

	return uc.repo.GetRandomSoldTicket(ctx, raffleID)
}

func (uc *TicketUseCase) GetRandomAvailableTicket(ctx context.Context, raffleID uint64) (*TicketWithOrganizerNumber, error) {
	return uc.repo.GetRandomAvailableTicket(ctx, raffleID)
}

func (uc *TicketUseCase) GenerateVoucher(ctx context.Context, ticketID uint64, userID uint64) ([]byte, error) {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	exists, err := uc.repo.IsTicketFromOrganizer(ctx, ticketID, organizer.ID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("ticket no encontrado para este organizador")
	}

	ticket, err := uc.repo.GetVoucherDataByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.FullName == "" {
		return nil, errors.New("el ticket no tiene un nombre completo asociado")
	}

	data := utils.TicketData{
		TicketID: ticket.Slug,
		Amount:   fmt.Sprintf("%d", ticket.Amount),
		DateTime: time.Now().Format("02 JAN, 2006 | 15:04"),
		FullName: ticket.FullName,
	}

	imgBytes, err := utils.GenerateTicketImage(data)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}

func (uc *TicketUseCase) Create(ctx context.Context, number uint64, participantName, participantPhone string, raffleID uint64) (*models.Ticket, error) {
	ticket, err := NewTicket(number, participantName, participantPhone, raffleID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (uc *TicketUseCase) Update(ctx context.Context, id uint64, participantName, participantPhone string, ticketStatusID uint64, userID uint64) (*models.Ticket, error) {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	exists, err := uc.repo.IsTicketFromOrganizer(ctx, id, organizer.ID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("ticket no encontrado para este organizador")
	}

	updateData, err := BuildTicketUpdateData(participantName, participantPhone, ticketStatusID)
	if err != nil {
		return nil, err
	}

	log.Printf("Update data for ticket ID %d: %v", id, updateData)

	if err := uc.repo.Update(ctx, id, updateData); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, id, organizer.ID)
}

func (uc *TicketUseCase) UpdateWinner(ctx context.Context, id uint64, raffleID uint64, userID uint64) error {

	if raffleID == 0 {
		return errors.New("el ID de la rifa es inválido")
	}

	if id == 0 {
		return errors.New("el ID del ticket es inválido")
	}

	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return errors.New("organizador no encontrado")
	}

	exists, err := uc.repo.IsTicketFromOrganizer(ctx, id, organizer.ID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("ticket no encontrado para este organizador")
	}

	if err := uc.repo.UpdateWinner(ctx, id, raffleID); err != nil {
		return err
	}

	return nil
}

func (uc *TicketUseCase) Delete(ctx context.Context, id uint64, userID uint64) error {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return errors.New("organizador no encontrado")
	}

	if _, err := uc.repo.GetByID(ctx, id, organizer.ID); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, id)
}
