package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser создает пользователя
// @Summary Создание пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.UserCreateRequest true "Данные пользователя"
// @Success 201 {object} entities.UserCreateResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /user/create [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Ограничение времени выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Парсинг данных из тела запроса
	userCreateRequest := &entities.UserCreateRequest{}
	if err := c.BodyParser(&userCreateRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Создание пользователя
	userCreateResponse, err := h.userService.CreateUser(ctx, userCreateRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(userCreateResponse)
}

// GetUserByID получение пользователя по ID
// @Summary получение пользователя по ID
// @Description получаем пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.UserGet
// @Failure 400 {object} entities.ErrorEntity
// @Failure 404 {object} entities.ErrorEntity
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// Ограничение времени выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Missing id"})
	}

	if userID < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(200).JSON(user)
}

// GetAllUsers получение всех пользователей
// @Summary получение всех пользователей
// @Description получение всех пользователей
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} []entities.UserGet
// @Failure 500 {object} entities.ErrorEntity
// @Router /user/all [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"users": users})
}
