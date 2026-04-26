package prize

import (
	"context"
	"spinLuck/internal/shared/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) PrizeRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) WithTransaction(fn func(repo PrizeRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &PostgresRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *PostgresRepository) GetAllByIdRaffleByOrganizerID(ctx context.Context, raffleID uint64, organizerID uint64) ([]models.Prize, error) {
	prizes := make([]models.Prize, 0)

	err := r.db.
		WithContext(ctx).
		Model(&models.Prize{}).
		Joins("JOIN raffle_prizes rp ON rp.prize_id = prizes.id").
		Where("rp.raffle_id = ? AND prizes.organizer_id = ?", raffleID, organizerID).
		Preload("Organizer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Find(&prizes).Error

	if err != nil {
		return nil, err
	}

	return prizes, nil
}

func (r *PostgresRepository) GetAllByIdRafflePublic(ctx context.Context, slug string) ([]models.Prize, error) {
	prizes := make([]models.Prize, 0)

	err := r.db.
		WithContext(ctx).
		Model(&models.Prize{}).
		Select("prizes.id, prizes.title, prizes.description, prizes.image_url").
		Joins("JOIN raffle_prizes rp ON rp.prize_id = prizes.id").
		Where("rp.raffle_id = (SELECT id FROM raffles WHERE slug = ? AND raffle_status_id = 1)", slug).
		Find(&prizes).Error

	if err != nil {
		return nil, err
	}

	return prizes, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint64) (*models.Prize, error) {
	var prize models.Prize
	if err := r.db.WithContext(ctx).First(&prize, id).Error; err != nil {
		return nil, err
	}
	return &prize, nil
}

func (r *PostgresRepository) ExistsByIDAndOrganizerID(ctx context.Context, id uint64, organizerID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Prize{}).Where("id = ? AND organizer_id = ?", id, organizerID).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (r *PostgresRepository) Create(ctx context.Context, prize *models.Prize) error {
	return r.db.WithContext(ctx).Create(prize).Error
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, data map[string]any) error {

	result := r.db.WithContext(ctx).
		Model(&models.Prize{}).
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

func (r *PostgresRepository) AssignToRaffle(ctx context.Context, raffleID, prizeID uint64) error {
	rafflePrize := models.RafflePrize{
		RaffleID: raffleID,
		PrizeID:  prizeID,
	}
	return r.db.WithContext(ctx).Create(&rafflePrize).Error
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Delete(&models.Prize{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteFromRaffle(ctx context.Context, raffleID uint64, prizeID uint64) error {
	result := r.db.WithContext(ctx).
		Where("raffle_id = ? AND prize_id = ?", raffleID, prizeID).
		Delete(&models.RafflePrize{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
