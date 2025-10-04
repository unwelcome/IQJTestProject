package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/services"
	"github.com/unwelcome/iqjtest/pkg/utils"
	"time"
)

type CatPhotoHandler struct {
	catPhotoService *services.CatPhotoService
}

func NewCatPhotoHandler(catPhotoService *services.CatPhotoService) *CatPhotoHandler {
	return &CatPhotoHandler{catPhotoService: catPhotoService}
}

// AddCatPhotos
// @Summary Загрузить фото кота
// @Description Загружает фотографию для указанного кота (не более 20 файлов)
// @Tags cat-photo
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param files formData []file true "Файлы изображений"
// @Success 201 {object} entities.CatPhotoUploadResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.CatPhotoUploadResponse
// @Router /auth/cat/mw/{id}/photo/add [post]
func (h *CatPhotoHandler) AddCatPhotos(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Получаем файлы из multipart/formData
	files, err := utils.GetFilesFromFormData(c, "files", 20)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)

	// Загружаем фото
	catPhotoUploadResponse := h.catPhotoService.AddCatPhoto(ctx, catID, files)

	// Отправляем результат
	if catPhotoUploadResponse.UploadedCount > 0 {
		return c.Status(201).JSON(catPhotoUploadResponse)
	}
	return c.Status(500).JSON(catPhotoUploadResponse)
}

// GetCatPhotoByID
// @Summary Получение фото кота
// @Description Получение всей информации о фото кота
// @Tags cat-photo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param photoID path int true "Photo ID"
// @Success 200 {object} entities.CatPhoto
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/photo/{photoID} [get]
func (h *CatPhotoHandler) GetCatPhotoByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	photoID, err := getPhotoID(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catPhoto, err := h.catPhotoService.GetCatPhotoByID(ctx, photoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(catPhoto)
}

// SetCatPhotoPrimary
// @Summary Выбор основного фото кота
// @Description Выбор основного фото кота
// @Tags cat-photo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param photoID path int true "Photo ID"
// @Success 200 {object} entities.CatPhotoSetPrimaryResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/mw/{id}/photo/{photoID}/primary [patch]
func (h *CatPhotoHandler) SetCatPhotoPrimary(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	photoID, err := getPhotoID(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)
	res, err := h.catPhotoService.SetCatPhotoPrimary(ctx, catID, photoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(res)
}

// DeleteCatPhoto
// @Summary Удаление фото кота по ID
// @Description Удаление фото кота по ID
// @Tags cat-photo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param photoID path int true "Photo ID"
// @Success 201 {object} string
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/mw/{id}/photo/{photoID} [delete]
func (h *CatPhotoHandler) DeleteCatPhoto(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	photoID, err := getPhotoID(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)

	err = h.catPhotoService.DeleteCatPhoto(ctx, catID, photoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).SendString("successfully deleted photo")
}

func getPhotoID(c *fiber.Ctx) (int, error) {
	photoID, err := c.ParamsInt("photoID")
	if err != nil {
		return 0, fmt.Errorf("missing photo id")
	}
	if photoID < 1 {
		return 0, fmt.Errorf("invalid photo id")
	}
	return photoID, nil
}
