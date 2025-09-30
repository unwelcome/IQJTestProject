package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/unwelcome/iqjtest/internal/entities"
	"time"
)

type AuthRepository struct {
	redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) *AuthRepository {
	return &AuthRepository{redis: redis}
}

func (r *AuthRepository) AddToken(ctx context.Context, userID int, token string, tokenType string, expiresIn time.Duration) error {
	key := getTokenKey(userID, tokenType)

	// Добавляем токен в сет
	err := r.redis.SAdd(ctx, key, token).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to add token: %s", err.Error()))
	}

	// Устанавливаем TTL на весь сет
	err = r.redis.Expire(ctx, key, expiresIn).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to set TTL to tokens: %s", err.Error()))
	}

	return nil
}

func (r *AuthRepository) CheckExistsToken(ctx context.Context, userID int, token string, tokenType string) error {
	key := getTokenKey(userID, tokenType)

	exists, err := r.redis.SIsMember(ctx, key, token).Result()
	if err != nil {
		return errors.New(fmt.Sprintf("check exists token failed: %s", err.Error()))
	}
	if !exists {
		return errors.New("token does not exist")
	}
	return nil
}

func (r *AuthRepository) ReplaceToken(ctx context.Context, userID int, oldToken, newToken, tokenType string, expiresIn time.Duration) error {
	key := getTokenKey(userID, tokenType)

	// Создаем транзакцию для замены токенов
	_, err := r.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		// Удаляем старый токен
		pipe.SRem(ctx, key, oldToken)

		// Добавляем новый токен
		pipe.SAdd(ctx, key, newToken)

		// Обновляем TTL для всего сета
		pipe.Expire(ctx, key, expiresIn)

		return nil
	})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to replace tokens: %s", err.Error()))
	}

	return nil
}

func (r *AuthRepository) DeleteToken(ctx context.Context, userID int, token string, tokenType string) error {
	key := getTokenKey(userID, tokenType)

	err := r.redis.SRem(ctx, key, token).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete token: %s", err.Error()))
	}
	return nil
}

func (r *AuthRepository) DeleteAllTokens(ctx context.Context, userID int) error {
	// Удаляем все access токены
	key := getTokenKey(userID, entities.AccessTokenType)
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete all access tokens: %s", err.Error()))
	}

	// Удаляем все refresh токены
	key = getTokenKey(userID, entities.RefreshTokenType)
	err = r.redis.Del(ctx, key).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete all refresh tokens: %s", err.Error()))
	}

	return nil
}

func getTokenKey(userID int, tokenType string) string {
	if tokenType == entities.AccessTokenType {
		return fmt.Sprintf("user:%d:access_tokens", userID)
	} else {
		return fmt.Sprintf("user:%d:refresh_tokens", userID)
	}
}
