package prize

import (
	"context"
	"errors"
	"log"
	"spinLuck/internal/features/organizer"
	"spinLuck/internal/features/storage"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"
)

type PrizeUseCase struct {
	repo  PrizeRepository
	ucOrg organizer.OrganizerService
	ucStg storage.StorageServiceInterface
}

func NewPrizeUseCase(repo PrizeRepository, ucOrg organizer.OrganizerService, ucStg storage.StorageServiceInterface) PrizeService {
	return &PrizeUseCase{repo: repo, ucOrg: ucOrg, ucStg: ucStg}
}

func (uc *PrizeUseCase) GetAllByIdRaffleByOrganizerID(ctx context.Context, userID uint64, raffleID uint64) ([]models.Prize, error) {
	return uc.repo.GetAllByIdRaffleByOrganizerID(ctx, raffleID, userID)
}

func (uc *PrizeUseCase) GetAllByIdRafflePublic(ctx context.Context, slug string) ([]models.Prize, error) {
	return uc.repo.GetAllByIdRafflePublic(ctx, slug)
}

func (uc *PrizeUseCase) GetByID(ctx context.Context, id uint64) (*models.Prize, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *PrizeUseCase) Create(ctx context.Context, userID uint64, title string, description string, imageURL string, raffleID uint64) (*models.Prize, error) {

	organizer, err := uc.ucOrg.GetBasicInfoByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	prize, err := NewPrize(title, description, imageURL, organizer.ID)
	if err != nil {
		return nil, err
	}

	err = uc.repo.WithTransaction(func(repo PrizeRepository) error {
		if err := repo.Create(ctx, prize); err != nil {
			return err
		}

		if raffleID > 0 {
			if err := repo.AssignToRaffle(ctx, raffleID, prize.ID); err != nil {
				return err
			}
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

	return prize, nil
}

func (uc *PrizeUseCase) Update(ctx context.Context, userID uint64, id uint64, title, description, imageURL string) (*models.Prize, error) {
	organizer, err := uc.ucOrg.GetBasicInfoByUserID(userID)
	if err != nil {
		return nil, errors.New("organizador no encontrado")
	}

	exists, err := uc.repo.ExistsByIDAndOrganizerID(ctx, id, organizer.ID)
	if err != nil {
		log.Printf("Error al verificar existencia de premio para actualización: %v", err)
		return nil, err
	}

	if !exists {
		return nil, errors.New("no tienes permiso para actualizar este premio o el premio no existe")
	}

	updateData, err := BuildPrizeUpdateData(title, description, imageURL)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, id, updateData); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *PrizeUseCase) Delete(ctx context.Context, userID uint64, raffleID uint64, prizeID uint64) error {
	organizer, err := uc.ucOrg.GetByUserID(userID)
	if err != nil {
		log.Printf("Error al obtener organizador para eliminación de premio: %v", err)
		return errors.New("organizador no encontrado")
	}

	prize, err := uc.repo.GetByID(ctx, prizeID)
	if err != nil {
		log.Printf("Error al obtener premio para eliminación: %v", err)
		return err
	}

	if prize.OrganizerID != organizer.ID {
		return errors.New("no tienes permiso para eliminar este premio")
	}

	err = uc.repo.WithTransaction(func(repo PrizeRepository) error {
		if err := repo.DeleteFromRaffle(ctx, raffleID, prizeID); err != nil {
			log.Printf("Error al eliminar premio de la rifa: %v", err)
			return err
		}

		url := utils.ExtractFileIdFromURL(prize.ImageURL)
		if err := uc.ucStg.Delete(ctx, userID, url); err != nil {
			log.Printf("Error al eliminar la imagen después de fallo en eliminación de premio: %v", err)
			return errors.New("error al eliminar la imagen después de fallo en eliminación de premio")
		}

		if err := repo.Delete(ctx, prizeID); err != nil {
			log.Printf("Error al eliminar premio: %v", err)
			return err
		}

		return nil
	})

	return err
}
