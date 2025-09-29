package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/services"
	"time"
)

func AuthMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем заголовок авторизации
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
		}

		// Проверяем корректность заголовка
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		// Получаем токен из заголовка
		refreshToken := authHeader[7:]

		// Создаем контекст
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Валидируем токен
		userID, err := authService.ValidateRefreshToken(ctx, refreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
		}

		// Устанавливаем userID и refresh_token в контекст
		c.Locals("userID", userID)
		c.Locals("refreshToken", refreshToken)

		return c.Next()
	}
}
