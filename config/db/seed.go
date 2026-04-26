package db

import "spinLuck/internal/shared/models"

func SeedDatabase() error {
	var countData int64

	if err := DB.Model(&models.RaffleStatus{}).Count(&countData).Error; err != nil {
		return err
	}

	if countData == 0 {
		if err := DB.Create(&models.RaffleStatuses).Error; err != nil {
			return err
		}
	}

	if err := DB.Model(&models.TicketStatus{}).Count(&countData).Error; err != nil {
		return err
	}

	if countData == 0 {
		if err := DB.Create(&models.TicketStatuses).Error; err != nil {
			return err
		}
	}

	return nil
}
