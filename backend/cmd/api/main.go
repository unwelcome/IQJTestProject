package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unwelcome/iqjtest/database/postgresql"
	"github.com/unwelcome/iqjtest/internal/config"
	"github.com/unwelcome/iqjtest/internal/dependency_injection"
	"github.com/unwelcome/iqjtest/internal/routes"
)

func main() {
	// Инициализация логгера
	zerolog.TimeFieldFormat = "02.01.2006 15:04:05.000"
	logger := log.With().Logger()

	// Инициализация конфига
	cfg := config.LoadConfig(logger)

	// Подключение к postgresql
	postgres := postgresql.Connect(cfg, logger)

	// Инициализация fiber
	app := fiber.New()

	// Создание контейнера с dependency injection
	container := dependency_injection.NewContainer(postgres.DB, logger)

	// Инициализация роутов
	routes.SetupRoutes(app, container)

	// Запуск приложения
	if err := app.Listen(":8080"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
