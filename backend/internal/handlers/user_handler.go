package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
	"time"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Ограничиваем время выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Парсим данные из тела запроса
	userCreateRequest := &entities.UserCreateRequest{}
	if err := c.BodyParser(&userCreateRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Создание пользователя
	userCreateResponse, err := h.userService.CreateUser(ctx, userCreateRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(userCreateResponse)
}
