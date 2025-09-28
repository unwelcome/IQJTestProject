package repositories

import (
	"context"
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

func (r *AuthRepository) SaveToken(ctx context.Context, token *entities.RefreshToken, userID int, tokenLifetime time.Duration) error {
	key := fmt.Sprintf("user-token-%d", userID)
	if err := r.redis.Set(ctx, key, token, tokenLifetime).Err(); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetToken(ctx context.Context, userID int) (*entities.RefreshToken, error) {
	key := fmt.Sprintf("user-token-%d", userID)
	token, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &entities.RefreshToken{Token: token}, nil
}

func (r *AuthRepository) DeleteToken(ctx context.Context, userID int) error {
	key := fmt.Sprintf("user-token-%d", userID)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
