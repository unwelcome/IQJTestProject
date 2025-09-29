package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
	"time"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register
// @Summary Создание пользователя
// @Description Создает нового пользователя в системе и возвращает access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entities.UserCreateRequest true "Данные пользователя"
// @Success 201 {object} entities.TokenPair
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Ограничение времени выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Парсинг данных из тела запроса
	userCreateRequest := &entities.UserCreateRequest{}
	if err := c.BodyParser(&userCreateRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Регистрация пользователя и получение токенов
	tokenPair, err := h.authService.RegistrationUser(ctx, userCreateRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(tokenPair)
}

// Login
// @Summary Вход в аккаунт пользователя
// @Description Вход в аккаунт пользователя, возвращает access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entities.UserLoginRequest true "Данные пользователя"
// @Success 200 {object} entities.TokenPair
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Ограничение времени выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userLoginRequest := &entities.UserLoginRequest{}
	if err := c.BodyParser(&userLoginRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	tokenPair, err := h.authService.LoginUser(ctx, userLoginRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(tokenPair)
}

// Refresh
// @Summary Обновление токенов
// @Description Обновляет access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} entities.TokenPair
// @Failure 400 {object} entities.ErrorEntity
// @Failure 401 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/refresh [get]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получаем userID и refresh токен из locals
	userID := c.Locals("userID").(int)
	refreshToken := c.Locals("refreshToken").(string)

	tokenPair, err := h.authService.RefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(tokenPair)
}

// Logout
// @Summary Удаление токена
// @Description Удаляет refresh токен
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Failure 401 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/logout [delete]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := c.Locals("userID").(int)
	refreshToken := c.Locals("refreshToken").(string)

	err := h.authService.DeleteRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).SendString("Successfully logged out")
}
