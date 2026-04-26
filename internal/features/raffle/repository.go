package raffle

import (
	"context"
	"spinLuck/internal/shared/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) RaffleRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAllStatus(ctx context.Context) ([]models.RaffleStatus, error) {
	var statuses []models.RaffleStatus
	if err := r.db.WithContext(ctx).Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]models.Raffle, error) {
	var raffles []models.Raffle
	if err := r.db.WithContext(ctx).Find(&raffles).Error; err != nil {
		return nil, err
	}
	return raffles, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Raffle, error) {
	var raffle models.Raffle
	if err := r.db.WithContext(ctx).First(&raffle, id).Error; err != nil {
		return nil, err
	}
	return &raffle, nil
}

func (r *PostgresRepository) GetAllInfoGeneric(ctx context.Context, organizerID uint64) ([]RaffleInfoGeneric, error) {
	var raffles []RaffleInfoGeneric
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			r.id,
			r.title,
			r.slug,
			r.quantity_tickets,
			rs.name AS status,
			r.image_url,

			COALESCE(SUM(CASE WHEN t.ticket_status_id = 1 THEN r.price ELSE 0 END), 0) AS total_amount,
			COALESCE(SUM(CASE WHEN t.ticket_status_id = 1 THEN 1 ELSE 0 END), 0) AS total_sold,
			CASE
				WHEN r.quantity_tickets = 0 THEN 0
				ELSE COALESCE(SUM(CASE WHEN t.ticket_status_id = 1 THEN 1 ELSE 0 END), 0) * 100.0 / r.quantity_tickets
			END AS progress
		FROM raffles r
		JOIN raffle_statuses rs ON r.raffle_status_id = rs.id
		LEFT JOIN tickets t ON r.id = t.raffle_id
		WHERE r.organizer_id = ?
		GROUP BY r.id, rs.name, r.created_at, r.price, r.quantity_tickets
		ORDER BY r.created_at DESC
	`, organizerID).Scan(&raffles).Error

	if err != nil {
		return nil, err
	}

	return raffles, nil
}

func (r *PostgresRepository) GetAllRecentInfoGeneric(ctx context.Context, organizerID uint64) ([]RaffleInfoGeneric, error) {
	var raffles []RaffleInfoGeneric
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			r.id,
			r.title,
			r.slug,
			r.quantity_tickets,
			rs.name AS status,
			r.image_url,

			COALESCE(SUM(CASE WHEN t.ticket_status_id = 1 THEN r.price ELSE 0 END), 0) AS total_amount,
			COALESCE(SUM(CASE WHEN t.ticket_status_id = 1 THEN 1 ELSE 0 END), 0) AS total_sold,
			CASE
				WHEN r.quantity_tickets = 0 THEN 0
				ELSE COALESCE(SUM(CASE WHEN t.ticket_status_id = 1 THEN 1 ELSE 0 END), 0) * 100.0 / r.quantity_tickets
			END AS progress
		FROM raffles r
		JOIN raffle_statuses rs ON r.raffle_status_id = rs.id
		LEFT JOIN tickets t ON r.id = t.raffle_id
		WHERE r.organizer_id = ?
		GROUP BY r.id, rs.name, r.created_at, r.price, r.quantity_tickets
		ORDER BY r.created_at DESC
		LIMIT 3
	`, organizerID).Scan(&raffles).Error

	if err != nil {
		return nil, err
	}

	return raffles, nil
}

func (r *PostgresRepository) GetByOrganizerID(ctx context.Context, organizerID uint64) ([]models.Raffle, error) {
	var raffles []models.Raffle
	if err := r.db.WithContext(ctx).Where("organizer_id = ?", organizerID).Find(&raffles).Error; err != nil {
		return nil, err
	}
	return raffles, nil
}

func (r *PostgresRepository) GetInfoBasicByID(ctx context.Context, id uint64, organizerID uint64) (*RaffleBasicInfo, error) {
	var info RaffleBasicInfo
	err := r.db.WithContext(ctx).
		Model(&models.Raffle{}).
		Select("id, organizer_id, image_url, quantity_tickets").
		Where("id = ? AND organizer_id = ?", id, organizerID).
		Take(&info).Error

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*models.Raffle, error) {
	var count int64
	var raffle models.Raffle

	if err := r.db.WithContext(ctx).
		Where("slug = ? AND raffle_status_id = ?", slug, 1).
		First(&raffle).Error; err != nil {
		return nil, err
	}

	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Where("raffle_id = ? AND ticket_status_id = ?", raffle.ID, 1).
		Count(&count).Error
	if err != nil {
		return nil, err
	}

	raffle.TicketsAvailable = raffle.QuantityTickets - uint64(count)

	return &raffle, nil
}

func (r *PostgresRepository) ExistsByIDAndOrganizerID(ctx context.Context, id uint64, organizerID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Raffle{}).Where("id = ? AND organizer_id = ?", id, organizerID).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (r *PostgresRepository) WithTransaction(fn func(repo RaffleRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &PostgresRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *PostgresRepository) Create(ctx context.Context, raffle *models.Raffle) error {
	return r.db.WithContext(ctx).Create(raffle).Error
}

func (r *PostgresRepository) CreateTickets(ctx context.Context, tickets []models.Ticket) error {
	return r.db.WithContext(ctx).Create(&tickets).Error
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data map[string]any) error {
	result := r.db.WithContext(ctx).
		Model(&models.Raffle{}).
		Where("id = ?", id).
		Updates(data)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresRepository) GetByIDForUpdate(ctx context.Context, id uint64) (*models.Raffle, error) {
	var raffle models.Raffle

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&raffle).Error
	if err != nil {
		return nil, err
	}

	return &raffle, nil
}

func (r *PostgresRepository) DeleteAvailableTicketsFromNumber(ctx context.Context, raffleID uint64, fromNumber uint64) error {
	return r.db.WithContext(ctx).
		Where("raffle_id = ? AND number >= ? AND ticket_status_id = ?", raffleID, fromNumber, 2).
		Delete(&models.Ticket{}).Error
}

func (r *PostgresRepository) CountNonAvailableTicketsFromNumber(ctx context.Context, raffleID uint64, fromNumber uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Where("raffle_id = ? AND number >= ? AND ticket_status_id <> ?", raffleID, fromNumber, 2).
		Count(&count).Error
	return count, err
}
