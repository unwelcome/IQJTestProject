package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/services"
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
		accessToken := authHeader[7:]

		// Создаем контекст
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Валидируем токен
		userID, err := authService.ValidateAccessToken(ctx, accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Устанавливаем userID в контекст
		c.Locals("userID", userID)

		return c.Next()
	}
}
