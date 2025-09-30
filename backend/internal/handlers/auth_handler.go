package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
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
// @Success 201 {object} entities.AuthResponse
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
	authResponse, err := h.authService.RegistrationUser(ctx, userCreateRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(authResponse)
}

// Login
// @Summary Вход в аккаунт пользователя
// @Description Вход в аккаунт пользователя, возвращает access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entities.UserLoginRequest true "Данные пользователя"
// @Success 200 {object} entities.AuthResponse
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

	authResponse, err := h.authService.LoginUser(ctx, userLoginRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(authResponse)
}

// Refresh
// @Summary Обновление токенов
// @Description Обновляет access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param token body entities.RefreshTokenRequest true "Refresh токен"
// @Success 201 {object} entities.TokenPair
// @Failure 400 {object} entities.ErrorEntity
// @Failure 401 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получаем refresh токен из тела
	refreshTokenRequest := &entities.RefreshTokenRequest{}
	if err := c.BodyParser(&refreshTokenRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	tokenPair, err := h.authService.RefreshToken(ctx, refreshTokenRequest.RefreshToken)
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
// @Param token body entities.LogoutTokenRequest true "Refresh токен"
// @Success 200 {object} string
// @Failure 401 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/logout [delete]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := c.Locals("userID").(int)

	// Получаем refresh токен из тела
	logoutTokenRequest := &entities.LogoutTokenRequest{}
	if err := c.BodyParser(&logoutTokenRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	err := h.authService.DeleteRefreshToken(ctx, userID, logoutTokenRequest.RefreshToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).SendString("Successfully logged out")
}

// DeleteUser
// @Summary Удаление пользователя
// @Description Удаляет пользователя из системы и отзывает все refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Failure 401 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /auth/user/delete [delete]
func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := c.Locals("userID").(int)

	err := h.authService.DeleteUser(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).SendString("Successfully deleted user")
}
