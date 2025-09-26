package routes

import (
	"github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/dependency_injection"
)

func SetupRoutes(app *fiber.App, container *dependency_injection.Container) {
	// Логирование всех запросов
	app.Use(container.LoggingMiddleware)

	// Группировка всех api роутов
	api := app.Group("/api")

	// Инициализация swagger
	// swag init -o ./api/docs --dir ./cmd/api,./internal/entities,./internal/handlers
	api.Get("/swagger/*", swagger.HandlerDefault)

	// Health запрос
	api.Get("/health", container.HealthHandler.Health)

	// User запросы
	api.Post("/user/create", container.UserHandler.CreateUser)
	//api.Get("/user/:id", container.UserHandler.GetUserById)
}
