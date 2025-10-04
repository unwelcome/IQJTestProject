package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
	"mime/multipart"
	"time"
)

type CatPhotoHandler struct {
	catPhotoService *services.CatPhotoService
}

func NewCatPhotoHandler(catPhotoService *services.CatPhotoService) *CatPhotoHandler {
	return &CatPhotoHandler{catPhotoService: catPhotoService}
}

// AddCatPhoto
// @Summary Загрузить фото кота
// @Description Загружает фотографию для указанного кота
// @Tags cat
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param file formData file true "Файл изображения"
// @Param is_primary formData bool false "Сделать главным фото"
// @Success 201 {object} entities.CatPhotoUploadResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id}/photo/add [post]
func (h *CatPhotoHandler) AddCatPhoto(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Парсим multipart/form-data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to parse multipart form"})
	}

	// Получаем файл
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "no file found in form"})
	}

	// Проверяем содержимое файла
	firstFile := files[0]
	if firstFile.Size == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "file is empty"})
	}

	// Проверяем тип файла
	if !isImageFile(firstFile) {
		return c.Status(400).JSON(fiber.Map{"error": "file must be an image file (jpg, png, webp)"})
	}

	// Открываем файл
	fileReader, err := firstFile.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to open file"})
	}
	defer fileReader.Close()

	// Получаем доп параметры
	isPrimary := c.FormValue("is_primary") == "true"
	catID := c.Locals("catID").(int)

	// Подготавливаем данные для передачи в сервис
	req := &entities.CatPhotoUploadRequest{
		File:      fileReader,
		FileSize:  firstFile.Size,
		FileName:  firstFile.Filename,
		MimeType:  firstFile.Header.Get("Content-Type"),
		IsPrimary: isPrimary,
	}

	// Загружаем фото
	res, err := h.catPhotoService.AddCatPhoto(ctx, catID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(res)
}

// SetCatPhotoPrimary
// @Summary Выбор основного фото кота
// @Description Выбор основного фото кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param photoID path int true "Photo ID"
// @Success 200 {object} entities.CatPhotoSetPrimaryResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id}/photo/{photoID}/primary [post]
func (h *CatPhotoHandler) SetCatPhotoPrimary(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	photoID, err := c.ParamsInt("photoID")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "missing photo id"})
	}

	if photoID < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid photo id"})
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
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param photoID path int true "Photo ID"
// @Success 201 {object} string
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id}/photo/{photoID} [delete]
func (h *CatPhotoHandler) DeleteCatPhoto(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	photoID, err := c.ParamsInt("photoID")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "missing photo id"})
	}

	if photoID < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid photo id"})
	}

	catID := c.Locals("catID").(int)

	err = h.catPhotoService.DeleteCatPhoto(ctx, catID, photoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).SendString("successfully deleted photo")
}

func isImageFile(fileHeader *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	return allowedTypes[contentType]
}
