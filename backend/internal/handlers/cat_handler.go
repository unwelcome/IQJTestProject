package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
)

type CatHandler struct {
	catService *services.CatService
}

func NewCatHandler(catService *services.CatService) *CatHandler {
	return &CatHandler{catService: catService}
}

// CreateCat
// @Summary Создание кота
// @Description Создает нового кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cat body entities.CatCreateRequest true "Данные кота"
// @Success 201 {object} entities.CatCreateResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/create [post]
func (h *CatHandler) CreateCat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createCatRequest := &entities.CatCreateRequest{}
	if err := c.BodyParser(&createCatRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID := c.Locals("userID").(int)

	createCatResponse, err := h.catService.CreateCat(ctx, userID, createCatRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(createCatResponse)
}

// GetCatByID
// @Summary Получение кота по ID
// @Description Получение кота по ID
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Success 201 {object} entities.Cat
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id} [get]
func (h *CatHandler) GetCatByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	catID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Missing id"})
	}

	if catID < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	cat, err := h.catService.GetCatByID(ctx, catID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(cat)
}

// GetAllCats
// @Summary Получение всех котов
// @Description Получение всех котов
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []entities.Cat
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/all [get]
func (h *CatHandler) GetAllCats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cats, err := h.catService.GetAllCats(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(cats)
}

// UpdateCatName
// @Summary Обновление клички кота
// @Description Обновление клички кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param cat body entities.CatUpdateNameRequest true "Данные кота"
// @Success 200 {object} entities.CatUpdateNameResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id}/name [patch]
func (h *CatHandler) UpdateCatName(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	catUpdateNameRequest := &entities.CatUpdateNameRequest{}
	if err := c.BodyParser(&catUpdateNameRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)

	catUpdateNameResponse, err := h.catService.UpdateCatName(ctx, catID, catUpdateNameRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(catUpdateNameResponse)
}

// UpdateCatAge
// @Summary Обновление возраста кота
// @Description Обновление возраста кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param cat body entities.CatUpdateAgeRequest true "Данные кота"
// @Success 200 {object} entities.CatUpdateAgeResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id}/age [patch]
func (h *CatHandler) UpdateCatAge(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	catUpdateAgeRequest := &entities.CatUpdateAgeRequest{}
	if err := c.BodyParser(&catUpdateAgeRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)

	catUpdateAgeResponse, err := h.catService.UpdateCatAge(ctx, catID, catUpdateAgeRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(catUpdateAgeResponse)
}

// UpdateCatDescription
// @Summary Обновление описания кота
// @Description Обновление описания кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param cat body entities.CatUpdateDescriptionRequest true "Данные кота"
// @Success 200 {object} entities.CatUpdateDescriptionResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id}/description [patch]
func (h *CatHandler) UpdateCatDescription(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	catUpdateDescriptionRequest := &entities.CatUpdateDescriptionRequest{}
	if err := c.BodyParser(&catUpdateDescriptionRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)

	catUpdateDescriptionResponse, err := h.catService.UpdateCatDescription(ctx, catID, catUpdateDescriptionRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(catUpdateDescriptionResponse)
}

// UpdateCat
// @Summary Обновление клички, возраста и описания кота
// @Description Обновление клички, возраста и описания кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param cat body entities.CatUpdateRequest true "Данные кота"
// @Success 200 {object} entities.CatUpdateResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id} [put]
func (h *CatHandler) UpdateCat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	catUpdateRequest := &entities.CatUpdateRequest{}
	if err := c.BodyParser(&catUpdateRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	catID := c.Locals("catID").(int)

	catUpdateResponse, err := h.catService.UpdateCat(ctx, catID, catUpdateRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(catUpdateResponse)
}

// DeleteCat
// @Summary Удаление кота
// @Description Удаление кота
// @Tags cat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Success 200 {object} string
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/cat/{id} [delete]
func (h *CatHandler) DeleteCat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	catID := c.Locals("catID").(int)

	err := h.catService.DeleteCat(ctx, catID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).SendString("Successfully deleted cat")
}
