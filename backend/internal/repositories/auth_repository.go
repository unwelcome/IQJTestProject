package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type AuthRepository struct {
	redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) *AuthRepository {
	return &AuthRepository{redis: redis}
}

func (r *AuthRepository) AddRefreshToken(ctx context.Context, userID int, refreshToken string, expiresIn time.Duration) error {
	key := getRefreshTokenKey(userID)

	// Добавляем токен в сет
	err := r.redis.SAdd(ctx, key, refreshToken).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to add refresh token: %s", err.Error()))
	}

	// Устанавливаем TTL на весь сет
	err = r.redis.Expire(ctx, key, expiresIn).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to set TTL to refresh tokens: %s", err.Error()))
	}

	return nil
}

func (r *AuthRepository) CheckExistsRefreshToken(ctx context.Context, userID int, refreshToken string) error {
	key := getRefreshTokenKey(userID)

	exists, err := r.redis.SIsMember(ctx, key, refreshToken).Result()
	if err != nil {
		return errors.New(fmt.Sprintf("check exists refresh token failed: %s", err.Error()))
	}
	if !exists {
		return errors.New("refresh token does not exist")
	}
	return nil
}

func (r *AuthRepository) ReplaceRefreshToken(ctx context.Context, userID int, oldRefreshToken, newRefreshToken string, expiresIn time.Duration) error {
	key := getRefreshTokenKey(userID)

	// Создаем транзакцию для замены токенов
	_, err := r.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		// Удаляем старый refresh токен
		pipe.SRem(ctx, key, oldRefreshToken)

		// Добавляем новый refresh токен
		pipe.SAdd(ctx, key, newRefreshToken)

		// Обновляем TTL для всего сета
		pipe.Expire(ctx, key, expiresIn)

		return nil
	})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to replace refresh tokens: %s", err.Error()))
	}

	return nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, userID int, refreshToken string) error {
	key := getRefreshTokenKey(userID)

	err := r.redis.SRem(ctx, key, refreshToken).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete refresh token: %s", err.Error()))
	}
	return nil
}

func (r *AuthRepository) DeleteAllRefreshTokens(ctx context.Context, userID int) error {
	key := getRefreshTokenKey(userID)

	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete all refresh tokens: %s", err.Error()))
	}
	return nil
}

func getRefreshTokenKey(userID int) string {
	return fmt.Sprintf("user:%d:refresh_tokens", userID)
}
