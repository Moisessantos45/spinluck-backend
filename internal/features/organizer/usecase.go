package organizer

import (
	"context"
	"errors"
	"spinLuck/internal/features/auth"
	"spinLuck/internal/shared/models"
)

type OrganizerUseCase struct {
	repo   OrganizerRepository
	rpAuth auth.AuthRepository
}

func NewOrganizerUseCase(repo OrganizerRepository, rpAuth auth.AuthRepository) OrganizerService {
	return &OrganizerUseCase{repo: repo, rpAuth: rpAuth}
}

func (uc *OrganizerUseCase) GetDashboardInfoDetailGeneric(ctx context.Context, userID uint64) (OrganizerDashboardMetrics, error) {
	organizer, err := uc.repo.GetByUserID(userID)
	if err != nil {
		return OrganizerDashboardMetrics{}, errors.New("organizador no encontrado")
	}

	return uc.repo.GetDashboardInfoDetailGeneric(ctx, organizer.ID)
}

func (uc *OrganizerUseCase) GetByUserID(userID uint64) (*models.Organizer, error) {
	if userID == 0 {
		return nil, errors.New("El ID del usuario no puede ser cero")
	}

	return uc.repo.GetByUserID(userID)
}

func (uc *OrganizerUseCase) GetByID(userID uint64) (*models.Organizer, error) {
	return uc.repo.GetByID(userID)
}

func (uc *OrganizerUseCase) GetBasicInfoByUserID(userID uint64) (*OrganizerBasicInfo, error) {
	if userID == 0 {
		return nil, errors.New("El ID del usuario no puede ser cero")
	}

	return uc.repo.GetBasicInfoByUserID(userID)
}

func (uc *OrganizerUseCase) Create(name string, phone string, userID uint64) (*models.Organizer, error) {
	organizer, err := NewOrganizer(name, phone, userID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(organizer); err != nil {
		return nil, err
	}

	if err := uc.rpAuth.ChangeCompletProfile(userID, true); err != nil {
		return nil, err
	}

	return organizer, nil
}

func (uc *OrganizerUseCase) Update(id uint64, name, phone string) (*models.Organizer, error) {
	organizer, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedData, err := BuildOrganizerUpdateData(name, phone)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Update(id, updatedData); err != nil {
		return nil, err
	}

	return organizer, nil
}
