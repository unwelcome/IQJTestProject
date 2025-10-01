package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/services"
)

func CatOwnershipMiddleware(catService *services.CatService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cansel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cansel()

		//Получаем catID
		catID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Missing cat ID"})
		}

		// Проверяем корректность catID
		if catID < 1 {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid cat ID"})
		}

		// Получаем userID
		userID := c.Locals("userID").(int)

		// Проверяем права на изменения кота
		hasRights, err := catService.CheckOwnershipRight(ctx, userID, catID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		} else if !hasRights {
			return c.Status(403).JSON(fiber.Map{"error": "not enough right for this operation"})
		}

		// Устанавливаем catID в Locals
		c.Locals("catID", catID)

		return c.Next()
	}
}
