package dependency_injection

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/unwelcome/iqjtest/internal/handlers"
	"github.com/unwelcome/iqjtest/internal/middlewares"
	"github.com/unwelcome/iqjtest/internal/repositories"
	"github.com/unwelcome/iqjtest/internal/services"
)

type Container struct {
	// Middleware
	LoggingMiddleware func(c *fiber.Ctx) error

	// Health
	HealthHandler *handlers.HealthHandler

	// User
	userRepository *repositories.UserRepository
	userService    *services.UserService
	UserHandler    *handlers.UserHandler
}

func NewContainer(postgres *sql.DB, logger zerolog.Logger) *Container {
	// Создание контейнера
	container := &Container{}

	// Инициализация middleware
	container.InitMiddlewares(logger)

	// Инициализация репозиториев
	container.InitRepositories(postgres)

	// Инициализация сервисов
	container.InitServices()

	// Инициализация хендлеров
	container.InitHandlers()

	return container
}

func (c *Container) InitMiddlewares(logger zerolog.Logger) {
	c.LoggingMiddleware = middlewares.LoggingRequest(logger)
}

func (c *Container) InitRepositories(postgres *sql.DB) {
	c.userRepository = repositories.NewUserRepository(postgres)
}

func (c *Container) InitServices() {
	c.userService = services.NewUserService(c.userRepository)
}

func (c *Container) InitHandlers() {
	c.HealthHandler = handlers.NewHealthHandler()
	c.UserHandler = handlers.NewUserHandler(c.userService)
}
