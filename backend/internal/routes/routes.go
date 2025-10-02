package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/dependency_injection"
)

func SetupRoutes(app *fiber.App, container *dependency_injection.Container) {
	// Логирование всех запросов
	app.Use(container.LoggingMiddleware)

	// Группировка всех api роутов
	api := app.Group("/api")

	// Проверка авторизации
	api.Use("/auth", container.AuthMiddleware)

	// Инициализация swagger
	// swag init -o ./api/docs --dir ./cmd/api,./internal/entities,./internal/handlers
	api.Get("/swagger/*", swagger.HandlerDefault)

	// Health запрос
	api.Get("/health", container.HealthHandler.Health)

	// Auth запросы
	api.Post("/register", container.AuthHandler.Register)
	api.Post("/login", container.AuthHandler.Login)
	api.Post("/refresh", container.AuthHandler.Refresh)
	api.Delete("/auth/logout", container.AuthHandler.Logout)
	api.Delete("/auth/user/delete", container.AuthHandler.DeleteUser)

	// User запросы
	api.Get("/auth/user/all", container.UserHandler.GetAllUsers)
	api.Get("/auth/user/:id", container.UserHandler.GetUserByID)
	api.Patch("/auth/user/password", container.UserHandler.UpdateUserPassword)

	// Cat запросы
	// Общие запросы
	api.Get("/auth/cat/all", container.CatHandler.GetAllCats)
	api.Get("/auth/cat/:id", container.CatHandler.GetCatByID)
	api.Post("/auth/cat/create", container.CatHandler.CreateCat)
	// Middleware проверки прав собственности
	api.Use("/auth/cat/:id", container.CatOwnershipMiddleware)
	// Запросы с middleware
	api.Post("/auth/cat/:id/photo", container.CatPhotoHandler.AddCatPhoto)
	api.Put("/auth/cat/:id", container.CatHandler.UpdateCat)
	api.Patch("/auth/cat/:id/name", container.CatHandler.UpdateCatName)
	api.Patch("/auth/cat/:id/age", container.CatHandler.UpdateCatAge)
	api.Patch("/auth/cat/:id/description", container.CatHandler.UpdateCatDescription)
	api.Delete("/auth/cat/:id", container.CatHandler.DeleteCat)
}
