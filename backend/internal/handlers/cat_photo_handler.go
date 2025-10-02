package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/services"
)

type CatPhotoHandler struct {
	catPhotoService *services.CatPhotoService
}

func NewCatPhotoHandler(catPhotoService *services.CatPhotoService) *CatPhotoHandler {
	return &CatPhotoHandler{catPhotoService: catPhotoService}
}

func (h *CatPhotoHandler) AddCatPhoto(c *fiber.Ctx) error {
	return nil
}
