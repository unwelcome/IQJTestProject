package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/unwelcome/iqjtest/database/minio"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	CacheHost     string
	CachePort     string
	CacheUser     string
	CachePassword string
	CacheDBName   string

	S3Host     string
	S3Port     string
	S3User     string
	S3Password string
	S3UseSSL   bool
	S3Buckets  map[string]string

	JWTSecret            string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

func LoadConfig(l zerolog.Logger) *Config {
	// Загружаем .env файл (игнорируем ошибку если файла нет)
	_ = godotenv.Load("../.env")

	// Создаем экземпляр конфига
	cfg := &Config{}

	// Инициализируем данные из env файла для основной DB
	cfg.DBHost = "localhost"
	cfg.DBPort = getEnv("POSTGRES_PORT", "5432")
	cfg.DBUser = getEnv("POSTGRES_USER", "postgres")
	cfg.DBPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	cfg.DBName = getEnv("POSTGRES_DB", "app")

	// Инициализируем данные из env файла для кеширования
	cfg.CacheHost = "localhost"
	cfg.CachePort = getEnv("REDIS_PORT", "6379")
	cfg.CacheUser = getEnv("REDIS_USER", "redis")
	cfg.CachePassword = getEnv("REDIS_PASSWORD", "redis")
	cfg.CacheDBName = getEnv("REDIS_DB", "0")

	// Инициализируем данные из env файла для s3 хранилища
	cfg.S3Host = "localhost"
	cfg.S3Port = getEnv("MINIO_PORT", "9000")
	cfg.S3User = getEnv("MINIO_USER", "minio")
	cfg.S3Password = getEnv("MINIO_PASSWORD", "minio")
	cfg.S3UseSSL = getEnvBool("MINIO_SSL", false)
	cfg.S3Buckets = map[string]string{
		"catPhotoBucket": "cat-photo-bucket",
	}

	// Для запуска через Docker
	if getEnv("IS_DOCKER", "") == "true" {
		cfg.DBHost = getEnv("POSTGRES_HOST", "postgres")
		cfg.CacheHost = getEnv("REDIS_HOST", "redis")
		cfg.S3Host = getEnv("MINIO_HOST", "minio")
	}

	// Инициализируем jwt секрет
	cfg.JWTSecret = getEnv("JWT_SECRET", "ultra-secret-key")
	cfg.AccessTokenLifetime = 5 * time.Minute
	cfg.RefreshTokenLifetime = 30 * 24 * time.Hour

	l.Trace().Str("DBHost", cfg.DBHost).Str("DBPort", cfg.DBPort).Msg("Postgres config")
	l.Trace().Str("CacheHost", cfg.CacheHost).Str("CachePort", cfg.CachePort).Msg("Redis config")
	l.Trace().Str("S3Host", cfg.S3Host).Str("S3Port", cfg.S3Port).Msg("Minio config")

	// Возвращаем экземпляр конфига
	return cfg
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) CacheConnString() string {
	//"redis://<user>:<pass>@localhost:6379/<db>"
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", c.CacheUser, c.CachePassword, c.CacheHost, c.CachePort, c.CacheDBName)
}

func (c *Config) S3ConnConfig() *minio.ConnectConfig {
	return &minio.ConnectConfig{
		Endpoint: fmt.Sprintf("%s:%s", c.S3Host, c.S3Port),
		Username: c.S3User,
		Password: c.S3Password,
		UseSSL:   c.S3UseSSL,
	}
}

func (c *Config) GetS3Buckets() []string {
	var buckets []string
	for _, bucketName := range c.S3Buckets {
		buckets = append(buckets, bucketName)
	}
	return buckets
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		result, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return result
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return result
	}
	return defaultValue
}
