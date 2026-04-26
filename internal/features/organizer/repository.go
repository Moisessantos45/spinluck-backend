package organizer

import (
	"context"
	"spinLuck/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) OrganizerRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetDashboardInfoDetailGeneric(
	ctx context.Context,
	organizerID uint64,
) (OrganizerDashboardMetrics, error) {
	var info OrganizerDashboardMetrics

	err := r.db.WithContext(ctx).Raw(`
		WITH active_raffles AS (
			SELECT
				r.id,
				r.quantity_tickets,
				r.date,
				r.price
			FROM raffles r
			WHERE r.organizer_id = ?
			  AND r.raffle_status_id = 1
		),
		sold_tickets_data AS (
			SELECT
				t.id,
				t.raffle_id,
				t.participant_phone,
				t.created_at,
				ar.price
			FROM tickets t
			JOIN active_raffles ar ON ar.id = t.raffle_id
			WHERE t.ticket_status_id = 1
		),
		totals AS (
			SELECT
				(SELECT COUNT(*) FROM active_raffles) AS total_raffles,
				COALESCE((SELECT SUM(quantity_tickets) FROM active_raffles), 0) AS total_tickets,
				COALESCE((SELECT COUNT(*) FROM sold_tickets_data), 0) AS total_sold,
				COALESCE((SELECT SUM(price) FROM sold_tickets_data), 0) AS total_revenue,
				COALESCE((SELECT AVG(price) FROM sold_tickets_data), 0) AS avg_ticket_price,
				COALESCE((SELECT COUNT(DISTINCT participant_phone) FROM sold_tickets_data), 0) AS total_unique_participants,
				COALESCE((SELECT COUNT(*) FROM sold_tickets_data WHERE created_at >= NOW() - INTERVAL '7 days'), 0) AS sold_last_7_days,
				COALESCE((SELECT COUNT(*) FROM active_raffles WHERE date <= CURRENT_DATE + INTERVAL '2 days'), 0) AS expiring_soon,
				COALESCE((
					SELECT COUNT(*)
					FROM active_raffles ar
					WHERE NOT EXISTS (SELECT 1 FROM raffle_prizes rp WHERE rp.raffle_id = ar.id)
				), 0) AS without_prizes,
				COALESCE((
					SELECT COUNT(*)
					FROM tickets t
					JOIN active_raffles ar ON ar.id = t.raffle_id
					WHERE t.ticket_status_id IN (3)
				), 0) AS stagnant_tickets
		)
		SELECT
			total_sold::bigint AS total_participants,
			total_raffles::bigint AS total_raffles,
			total_unique_participants::bigint AS total_unique_participants,
			total_revenue::float8 AS total_amount,
			avg_ticket_price::float8 AS average_ticket_price,
			CASE
				WHEN total_tickets = 0 THEN 0
				ELSE ROUND((total_sold * 100.0) / total_tickets, 2)
			END AS effective_sales_rate,
			expiring_soon::bigint AS active_raffles_expiring_soon,
			CASE
				WHEN sold_last_7_days = 0 THEN 0
				ELSE ROUND(((total_tickets - total_sold) * 1.0) / NULLIF(sold_last_7_days / 7.0, 0), 2)
			END AS estimated_depletion_days,
			stagnant_tickets::bigint AS stagnant_tickets,
			without_prizes::bigint AS active_raffles_without_prizes
		FROM totals
	`, organizerID).Scan(&info).Error

	if err != nil {
		return OrganizerDashboardMetrics{}, err
	}

	return info, nil
}

func (r *PostgresRepository) GetByUserID(userID uint64) (*models.Organizer, error) {
	var organizer models.Organizer

	if err := r.db.Model(&models.Organizer{}).Where("user_id = ?", userID).First(&organizer).Error; err != nil {
		return nil, err
	}

	return &organizer, nil
}

func (r *PostgresRepository) GetByID(userID uint64) (*models.Organizer, error) {
	var organizer models.Organizer

	if err := r.db.Model(&models.Organizer{}).Where("user_id = ?", userID).First(&organizer).Error; err != nil {
		return nil, err
	}

	return &organizer, nil
}

func (r *PostgresRepository) GetBasicInfoByUserID(userID uint64) (*OrganizerBasicInfo, error) {
	var info OrganizerBasicInfo

	if err := r.db.Model(&models.Organizer{}).Select("id, user_id").Where("user_id = ?", userID).Take(&info).Error; err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *PostgresRepository) Create(organizer *models.Organizer) error {
	return r.db.Model(&models.Organizer{}).Create(organizer).Error
}

func (r *PostgresRepository) Update(id uint64, data map[string]any) error {
	return r.db.Model(&models.Organizer{}).Where("id = ?", id).Updates(data).Error
}
