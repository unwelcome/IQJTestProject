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

func (s *CatPhotoService) AddCatPhoto(ctx context.Context, catID int, req *entities.CatPhotoUploadRequest) (*entities.CatPhotoUploadSuccess, error) {
	// Загружаем фото кота
	res, err := s.catPhotoRepository.AddCatPhoto(ctx, catID, req)
	if err != nil {
		return nil, fmt.Errorf("add cat photo error: %w", err)
	}

	return res, nil
}

func (s *CatPhotoService) DeleteCatPhoto(ctx context.Context, catID, photoID int) error {
	// Получаем информацию о фото
	catPhoto, err := s.catPhotoRepository.GetCatPhotoByID(ctx, photoID)
	if err != nil {
		return fmt.Errorf("get cat photo error: %w", err)
	}

	// Проверяем, что фото принадлежит коту
	if catPhoto.CatID != catID {
		return fmt.Errorf("photo %d doesn't belong to cat %d", photoID, catID)
	}

	// Удаляем фото
	err = s.catPhotoRepository.DeleteCatPhoto(ctx, photoID)
	if err != nil {
		return fmt.Errorf("delete cat photo error: %w", err)
	}
	return nil
}

func (s *CatPhotoService) SetCatPhotoPrimary(ctx context.Context, catID int, photoID int) (*entities.CatPhotoSetPrimaryResponse, error) {
	err := s.catPhotoRepository.SetCatPhotoPrimary(ctx, catID, photoID)
	if err != nil {
		return nil, fmt.Errorf("set cat photo primary error: %w", err)
	}
	return &entities.CatPhotoSetPrimaryResponse{ID: photoID}, nil
}

func (s *CatPhotoService) GetCatPhotoByID(ctx context.Context, photoID int) (*entities.CatPhoto, error) {
	catPhoto, err := s.catPhotoRepository.GetCatPhotoByID(ctx, photoID)
	if err != nil {
		return nil, fmt.Errorf("get cat photo by id error: %w", err)
	}
	return catPhoto, nil
}
