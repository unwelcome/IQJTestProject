package services

import (
	"context"
	"fmt"

	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/repositories"
)

type CatPhotoService struct {
	catPhotoRepository *repositories.CatPhotoRepository
}

func NewCatPhotoService(catPhotoRepository *repositories.CatPhotoRepository) *CatPhotoService {
	return &CatPhotoService{catPhotoRepository: catPhotoRepository}
}

func (s *CatPhotoService) AddCatPhoto(ctx context.Context, catID int, req *entities.CatPhotoUploadRequest) (*entities.CatPhotoUploadResponse, error) {
	res, err := s.catPhotoRepository.AddCatPhoto(ctx, catID, req)
	if err != nil {
		return nil, fmt.Errorf("add cat photo error: %w", err)
	}

	return res, nil
}

func (s *CatPhotoService) DeleteCatPhoto(ctx context.Context, photoID int) error {
	err := s.catPhotoRepository.DeleteCatPhoto(ctx, photoID)
	if err != nil {
		return fmt.Errorf("delete cat photo error: %w", err)
	}
	return nil
}

func (s *CatPhotoService) SetCatPhotoPrimary(ctx context.Context, catID int, photoID int) error {
	err := s.catPhotoRepository.SetCatPhotoPrimary(ctx, catID, photoID)
	if err != nil {
		return fmt.Errorf("set cat photo primary error: %w", err)
	}
	return nil
}

func (s *CatPhotoService) GetCatPhotoByID(ctx context.Context, photoID int) (*entities.CatPhoto, error) {
	catPhoto, err := s.catPhotoRepository.GetCatPhotoByID(ctx, photoID)
	if err != nil {
		return nil, fmt.Errorf("get cat photo by id error: %w", err)
	}
	return catPhoto, nil
}
