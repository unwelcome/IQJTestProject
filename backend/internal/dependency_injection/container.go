package dependency_injection

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/unwelcome/iqjtest/internal/config"
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

	//Auth
	authRepository *repositories.AuthRepository
	authService    *services.AuthService
	AuthHandler    *handlers.AuthHandler

	// User
	userRepository *repositories.UserRepository
	userService    *services.UserService
	UserHandler    *handlers.UserHandler
}

func NewContainer(postgres *sql.DB, redis *redis.Client, cfg *config.Config, logger zerolog.Logger) *Container {
	// Создание контейнера
	container := &Container{}

	// Инициализация middleware
	container.InitMiddlewares(logger)

	// Инициализация репозиториев
	container.InitRepositories(postgres, redis)

	// Инициализация сервисов
	container.InitServices(cfg)

	// Инициализация хендлеров
	container.InitHandlers()

	return container
}

func (c *Container) InitMiddlewares(logger zerolog.Logger) {
	c.LoggingMiddleware = middlewares.LoggingRequest(logger)
}

func (c *Container) InitRepositories(postgres *sql.DB, redis *redis.Client) {
	c.userRepository = repositories.NewUserRepository(postgres)
	c.authRepository = repositories.NewAuthRepository(redis)
}

func (c *Container) InitServices(cfg *config.Config) {
	c.userService = services.NewUserService(c.userRepository)
	c.authService = services.NewAuthService(c.userService, c.authRepository, cfg.JWTSecret, cfg.AccessTokenLifetime, cfg.RefreshTokenLifetime)
}

func (c *Container) InitHandlers() {
	c.HealthHandler = handlers.NewHealthHandler()
	c.UserHandler = handlers.NewUserHandler(c.userService)
	c.AuthHandler = handlers.NewAuthHandler(c.authService)
}
