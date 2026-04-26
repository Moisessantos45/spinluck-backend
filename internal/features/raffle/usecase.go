package raffle

import (
	"context"
	"errors"
	"log"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/features/prize"
	"spinLuck/internal/features/storage"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"
)

type RaffleUseCase struct {
	repo  RaffleRepository
	rpOrg organizer.OrganizerRepository
	ucStg storage.StorageServiceInterface
	ucPrz prize.PrizeService
}

func NewRaffleUseCase(repo RaffleRepository, repoOgz organizer.OrganizerRepository, stg storage.StorageServiceInterface, ucPrz prize.PrizeService) RaffleService {
	return &RaffleUseCase{repo: repo, rpOrg: repoOgz, ucStg: stg, ucPrz: ucPrz}
}

func (uc *RaffleUseCase) GetAllStatus(ctx context.Context) ([]models.RaffleStatus, error) {
	return uc.repo.GetAllStatus(ctx)
}

func (uc *RaffleUseCase) GetAll(ctx context.Context) ([]models.Raffle, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *RaffleUseCase) GetAllInfoGeneric(ctx context.Context, userID uint64) ([]RaffleInfoGeneric, error) {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	return uc.repo.GetAllInfoGeneric(ctx, organizer.ID)
}

func (uc *RaffleUseCase) GetAllRecentInfoGeneric(ctx context.Context, userID uint64) ([]RaffleInfoGeneric, error) {
	organizer, err := uc.rpOrg.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	return uc.repo.GetAllRecentInfoGeneric(ctx, organizer.ID)
}

func (uc *RaffleUseCase) GetByID(ctx context.Context, id uint64) (*models.Raffle, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *RaffleUseCase) GetByOrganizerID(ctx context.Context, organizerID uint64) ([]models.Raffle, error) {
	return uc.repo.GetByOrganizerID(ctx, organizerID)
}

func (uc *RaffleUseCase) GetBySlug(ctx context.Context, userID uint64, slug string) (*models.Raffle, error) {

	result, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if userID == 0 {
		result.OrganizerID = 0
		result.RaffleStatusID = 0
	}

	return result, nil
}

func (uc *RaffleUseCase) Create(
	ctx context.Context,
	userID uint64,
	title string,
	description string,
	price float64,
	imageURL string,
	date string,
	quantityTickets uint64,
	maxWinners uint64,
	raffleStatusID uint64,
) (*models.Raffle, error) {
	organizer, err := uc.rpOrg.GetBasicInfoByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	raffle, err := NewRaffle(title, description, price, imageURL, date, quantityTickets, maxWinners, organizer.ID, raffleStatusID)
	if err != nil {
		return nil, err
	}

	err = uc.repo.WithTransaction(func(repo RaffleRepository) error {
		if err := repo.Create(ctx, raffle); err != nil {
			return err
		}

		tickets := make([]models.Ticket, 0, quantityTickets)
		for i := uint64(1); i <= quantityTickets; i++ {
			tickets = append(tickets, models.Ticket{
				RaffleID: raffle.ID,
				Number:   i,
			})
		}

		if err := repo.CreateTickets(ctx, tickets); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		url := utils.ExtractFileIdFromURL(imageURL)
		if err := uc.ucStg.Delete(ctx, userID, url); err != nil {
			log.Printf("Error al eliminar la imagen después de fallo en creación de rifa: %v", err)
			return nil, errors.New("error al eliminar la imagen después de fallo en creación de rifa")
		}
		return nil, err
	}

	return raffle, nil
}

func (uc *RaffleUseCase) Update(
	ctx context.Context,
	userID uint64,
	id uint64,
	title string,
	description string,
	price float64,
	imageURL string,
	date string,
	quantityTickets uint64,
	maxWinners uint64,
	raffleStatusID uint64,
) (*models.Raffle, error) {
	organizer, err := uc.rpOrg.GetBasicInfoByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	updateData, err := BuildRaffleUpdateData(
		title,
		description,
		price,
		imageURL,
		date,
		quantityTickets,
		maxWinners,
		raffleStatusID,
	)
	if err != nil {
		return nil, err
	}

	err = uc.repo.WithTransaction(func(repo RaffleRepository) error {
		currentRaffle, err := repo.GetInfoBasicByID(ctx, id, organizer.ID)
		if err != nil {
			return err
		}

		oldQty := currentRaffle.Quantity
		newQty := quantityTickets

		if newQty < oldQty {
			count, err := repo.CountNonAvailableTicketsFromNumber(ctx, id, newQty+1)
			if err != nil {
				return err
			}

			if count > 0 {
				return errors.New("no se puede reducir la cantidad: existen boletos vendidos o reservados fuera del nuevo rango")
			}
		}

		if err := repo.Update(ctx, id, updateData); err != nil {
			return err
		}

		if newQty > oldQty {
			tickets := make([]models.Ticket, 0, newQty-oldQty)
			for i := oldQty + 1; i <= newQty; i++ {
				tickets = append(tickets, models.Ticket{
					RaffleID:       id,
					Number:         i,
					TicketStatusID: 2,
				})
			}

			if len(tickets) > 0 {
				if err := repo.CreateTickets(ctx, tickets); err != nil {
					return err
				}
			}
		}

		if newQty < oldQty {
			if err := repo.DeleteAvailableTicketsFromNumber(ctx, id, newQty+1); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedRaffle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedRaffle, nil
}
