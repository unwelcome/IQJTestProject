package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unwelcome/iqjtest/database/postgresql"
	"github.com/unwelcome/iqjtest/internal/config"
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

	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info().Msg("Request get")
		return c.SendString("Hello, World!")
	})

	app.Post("/", func(c *fiber.Ctx) error {
		logger.Info().Msg("Request post")
		if err := postgres.DB.Ping(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).SendString("Pong")
	})

	// Запуск приложения
	if err := app.Listen(":8080"); err != nil {
		logger.Fatal().Err(err)
	}
}
