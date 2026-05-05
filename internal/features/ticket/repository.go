package ticket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"spinLuck/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) TicketRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAllStatus(ctx context.Context) ([]models.TicketStatus, error) {
	var statuses []models.TicketStatus
	if err := r.db.WithContext(ctx).Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := r.db.WithContext(ctx).Preload("Raffle").Preload("TicketStatus").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *PostgresRepository) GetAllByRaffleIDToOrgaizer(ctx context.Context, raffleID uint64, organizerID uint64) ([]models.Ticket, error) {
	var tickets []models.Ticket

	if err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Select("tickets.id, tickets.number,tickets.winner, tickets.raffle_id, tickets.ticket_status_id").
		Joins("JOIN raffles ON raffles.id = tickets.raffle_id").
		Where("raffles.id = ? AND raffles.organizer_id = ?", raffleID, organizerID).
		Preload("Raffle", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "organizer_id").Where("id = ?", raffleID)
		}).
		Preload("TicketStatus", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	for i := range tickets {
		tickets[i].FormattedNumber = fmt.Sprintf("%03d", tickets[i].Number)
	}

	return tickets, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64, organizerID uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.
		WithContext(ctx).
		Model(&models.Ticket{}).
		Joins("JOIN raffles ON raffles.id = tickets.raffle_id").
		Where("tickets.id = ? AND raffles.organizer_id = ?", id, organizerID).
		Preload("Raffle", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "organizer_id")
		}).
		Preload("TicketStatus", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(&ticket, id).Error; err != nil {
		return nil, err
	}

	ticket.FormattedNumber = fmt.Sprintf("%03d", ticket.Number)

	return &ticket, nil
}

func (r *PostgresRepository) IsTicketFromOrganizer(ctx context.Context, ticketID uint64, organizerID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Joins("JOIN raffles ON raffles.id = tickets.raffle_id").
		Where("tickets.id = ? AND raffles.organizer_id = ?", ticketID, organizerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresRepository) GetVoucherDataByID(ctx context.Context, id uint64) (*TicketVoucherData, error) {
	var data TicketVoucherData
	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Select("tickets.id as ticket_id, raffles.price as amount, tickets.participant_name as full_name,tickets.number as number, raffles.slug as slug").
		Joins("JOIN raffles ON raffles.id = tickets.raffle_id").
		Where("tickets.id = ?", id).
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	data.FormattedNumber = fmt.Sprintf("%03d", data.Number)

	return &data, nil
}

func (r *PostgresRepository) GetRandomSoldTicket(ctx context.Context, raffleID uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	var maxReached bool

	err := r.db.WithContext(ctx).
		Raw(`
        SELECT
            CASE
                WHEN COUNT(t.id) >= r.max_winners THEN true
                ELSE false
            END as max_winners_reached
        FROM raffles r
        LEFT JOIN tickets t ON r.id = t.raffle_id AND t.winner = true
        WHERE r.id = ? AND r.raffle_status_id = 1 ? AND r.date >= CURRENT_DATE
        GROUP BY r.id, r.max_winners
    `, raffleID).Scan(&maxReached).Error

	if err != nil {
		return nil, err
	}

	if maxReached {
		return nil, errors.New("se ha alcanzado el número máximo de ganadores para esta rifa")
	}

	err = r.db.WithContext(ctx).
		Raw(`
			SELECT id, number, participant_name, participant_phone
			FROM tickets
			WHERE raffle_id = ?
			  AND ticket_status_id = 1
			ORDER BY RANDOM()
			LIMIT 1;
			`, raffleID).
		Scan(&ticket).Error

	if err != nil {
		return nil, err
	}

	ticket.FormattedNumber = fmt.Sprintf("%03d", ticket.Number)

	return &ticket, nil
}

func (r *PostgresRepository) GetRandomAvailableTicket(ctx context.Context, raffleID uint64) (*TicketWithOrganizerNumber, error) {
	var ticket TicketWithOrganizerNumber

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				tickets.id,
				tickets.number,
				rf.slug,
				org.phone AS organizer_number
			FROM tickets
			INNER JOIN raffles rf ON rf.id = tickets.raffle_id
			INNER JOIN organizers org ON org.id = rf.organizer_id
			WHERE tickets.raffle_id = ? AND rf.raffle_status_id = 1
			  AND tickets.ticket_status_id = 2 ? AND rf.date >= CURRENT_DATE
			ORDER BY RANDOM()
			LIMIT 1
		`, raffleID).
		Scan(&ticket).Error

	if err != nil {
		return nil, err
	}

	ticket.FormattedNumber = fmt.Sprintf("%03d", ticket.Number)

	return &ticket, nil
}

func (r *PostgresRepository) FindByNumberAndRaffleID(ctx context.Context, number uint64, raffleID uint64) (*TicketWithOrganizerNumber, error) {
	var ticket TicketWithOrganizerNumber

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				tickets.id,
				tickets.number,
				rf.slug,
				org.phone AS organizer_number
			FROM tickets
			INNER JOIN raffles rf ON rf.id = tickets.raffle_id
			INNER JOIN organizers org ON org.id = rf.organizer_id
			WHERE tickets.number = ? AND tickets.raffle_id = ? AND rf.raffle_status_id = 1
			  AND rf.date >= CURRENT_DATE AND tickets.ticket_status_id = 2
		`, number, raffleID).
		Scan(&ticket).Error

	if err != nil {
		return nil, err
	}

	log.Printf("Número de ticket encontrado: %d", ticket.Number)
	ticket.FormattedNumber = fmt.Sprintf("%03d", ticket.Number)

	return &ticket, nil
}

func (r *PostgresRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Model(&models.Ticket{}).Create(ticket).Error
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", id).Updates(data).Error
}

func (r *PostgresRepository) UpdateWinner(ctx context.Context, ticketID uint64, raffleID uint64) error {
	result := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Where("id = ? AND raffle_id = ? AND winner = false", ticketID, raffleID).
		Update("winner", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("ticket no encontrado, no pertenece a la rifa o ya era ganador")
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Ticket{}, id).Error
}
