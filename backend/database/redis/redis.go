package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/unwelcome/iqjtest/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config, l zerolog.Logger) *redis.Client {
	opt, err := redis.ParseURL(cfg.CacheConnString())
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	if err = rdb.Ping(ctx).Err(); err != nil {
		l.Fatal().Err(err).Msg("failed to connect to redis server")
		return nil
	}

	l.Trace().Msg("Successfully connected to Redis!")
	return rdb
}
