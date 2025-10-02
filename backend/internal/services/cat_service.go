package services

import (
	"context"
	"fmt"

	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/repositories"
)

type CatService struct {
	catRepository *repositories.CatRepository
}

func NewCatService(catRepository *repositories.CatRepository) *CatService {
	return &CatService{catRepository: catRepository}
}

func (s *CatService) CreateCat(ctx context.Context, userID int, catCreateRequest *entities.CatCreateRequest) (*entities.CatCreateResponse, error) {
	cat := &entities.Cat{
		Name:        catCreateRequest.Name,
		Age:         catCreateRequest.Age,
		Description: catCreateRequest.Description,
	}

	err := s.catRepository.CreateCat(ctx, userID, cat)
	if err != nil {
		return nil, fmt.Errorf("create cat error: %s", err.Error())
	}
	return &entities.CatCreateResponse{ID: cat.ID}, nil
}

func (s *CatService) GetCatByID(ctx context.Context, catID int) (*entities.CatWithPhotos, error) {
	cat, err := s.catRepository.GetCatByID(ctx, catID)
	if err != nil {
		return nil, fmt.Errorf("get cat by id error: %s", err.Error())
	}
	return cat, nil
}

func (s *CatService) GetAllCats(ctx context.Context) ([]*entities.CatWithPrimePhoto, error) {
	cats, err := s.catRepository.GetAllCats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all cats error: %s", err.Error())
	}
	return cats, nil
}

func (s *CatService) UpdateCatName(ctx context.Context, catID int, catUpdateNameRequest *entities.CatUpdateNameRequest) (*entities.CatUpdateNameResponse, error) {
	err := s.catRepository.UpdateCatName(ctx, catID, catUpdateNameRequest.Name)
	if err != nil {
		return nil, fmt.Errorf("update cat name error: %s", err.Error())
	}
	return &entities.CatUpdateNameResponse{ID: catID}, nil
}

func (s *CatService) UpdateCatAge(ctx context.Context, catID int, catUpdateAgeRequest *entities.CatUpdateAgeRequest) (*entities.CatUpdateAgeResponse, error) {
	err := s.catRepository.UpdateCatAge(ctx, catID, catUpdateAgeRequest.Age)
	if err != nil {
		return nil, fmt.Errorf("update cat age error: %s", err.Error())
	}
	return &entities.CatUpdateAgeResponse{ID: catID}, nil
}

func (s *CatService) UpdateCatDescription(ctx context.Context, catID int, catUpdateDescriptionRequest *entities.CatUpdateDescriptionRequest) (*entities.CatUpdateDescriptionResponse, error) {
	err := s.catRepository.UpdateCatDescription(ctx, catID, catUpdateDescriptionRequest.Description)
	if err != nil {
		return nil, fmt.Errorf("update cat description error: %s", err.Error())
	}
	return &entities.CatUpdateDescriptionResponse{ID: catID}, nil
}

func (s *CatService) UpdateCat(ctx context.Context, catID int, catUpdateRequest *entities.CatUpdateRequest) (*entities.CatUpdateResponse, error) {
	err := s.catRepository.UpdateCat(ctx, catID, catUpdateRequest)
	if err != nil {
		return nil, fmt.Errorf("update cat error: %s", err.Error())
	}
	return &entities.CatUpdateResponse{
		ID:          catID,
		Name:        catUpdateRequest.Name,
		Age:         catUpdateRequest.Age,
		Description: catUpdateRequest.Description,
	}, nil
}

func (s *CatService) DeleteCat(ctx context.Context, catID int) error {
	err := s.catRepository.DeleteCat(ctx, catID)
	if err != nil {
		return fmt.Errorf("delete cat error: %s", err.Error())
	}
	return nil
}

func (s *CatService) CheckOwnershipRight(ctx context.Context, userID, catID int) (bool, error) {
	cat, err := s.catRepository.GetCatByID(ctx, catID)
	if err != nil {
		return false, fmt.Errorf("get cat info error: %s", err.Error())
	}

	return cat.CreatedBy == userID, nil
}
