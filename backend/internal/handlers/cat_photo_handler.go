package handlers

import (
	"context"
	"fmt"
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

	// Парсим multipart/form-data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to parse multipart form"})
	}

	// Получаем файлы
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "no file found in form"})
	}

	if len(files) > 20 {
		return c.Status(400).JSON(fiber.Map{"error": "too many files in form"})
	}

	catID := c.Locals("catID").(int)
	var uploadedPhotos []*entities.CatPhotoUploadSuccess
	var errors []*entities.CatPhotoUploadError

	// Проходимся по каждому фото
	for _, file := range files {

		// Проверяем содержимое файла
		if file.Size == 0 {
			errors = append(errors, &entities.CatPhotoUploadError{
				FileName: file.Filename,
				Error:    "file is empty",
			})
			continue
		}

		// Проверяем размер файла (макс. 50МБ)
		if file.Size > 50*1024*1024 {
			errors = append(errors, &entities.CatPhotoUploadError{
				FileName: file.Filename,
				Error:    "file is too large",
			})
			continue
		}

		// Проверяем тип файла
		if !isImageFile(file) {
			errors = append(errors, &entities.CatPhotoUploadError{
				FileName: file.Filename,
				Error:    "file must be an image file (jpg, png, webp)",
			})
			continue
		}

		// Открываем файл
		fileReader, err := file.Open()
		if err != nil {
			errors = append(errors, &entities.CatPhotoUploadError{
				FileName: file.Filename,
				Error:    "failed to open file",
			})
			continue
		}
		defer fileReader.Close()

		// Подготавливаем данные для передачи в сервис
		req := &entities.CatPhotoUploadRequest{
			File:     fileReader,
			FileSize: file.Size,
			FileName: file.Filename,
			MimeType: file.Header.Get("Content-Type"),
		}

		// Загружаем фото
		success, err := h.catPhotoService.AddCatPhoto(ctx, catID, req)
		if err != nil {
			errors = append(errors, &entities.CatPhotoUploadError{
				FileName: file.Filename,
				Error:    "failed to save file",
			})
			continue
		}

		// Добавляем данные фото в массив
		uploadedPhotos = append(uploadedPhotos, success)
	}

	// Подготавливаем тело ответа
	response := &entities.CatPhotoUploadResponse{
		Message:        fmt.Sprintf("Uploaded %d out of %d files", len(uploadedPhotos), len(files)),
		UploadedCount:  len(uploadedPhotos),
		FailedCount:    len(errors),
		UploadedPhotos: uploadedPhotos,
		Errors:         errors,
	}

	// Отправляем результат
	if len(uploadedPhotos) > 0 {
		return c.Status(201).JSON(response)
	}
	return c.Status(500).JSON(response)
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
