package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/dependency_injection"
)

func SetupRoutes(app *fiber.App, container *dependency_injection.Container) {
	// Логирование всех запросов
	app.Use(container.LoggingMiddleware)

	// Группировка всех api роутов
	api := app.Group("/api")

	// Health запрос
	api.Get("/health", container.HealthHandler.Health)

	// User запросы
	api.Post("/user/create", container.UserHandler.CreateUser)
	//api.Get("/user/:id", container.UserHandler.GetUserById)
}
