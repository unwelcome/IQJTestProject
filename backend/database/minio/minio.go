package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

type ConnectConfig struct {
	Endpoint string
	Username string
	Password string
	UseSSL   bool
}

func InitMinio(cfg *ConnectConfig, l zerolog.Logger, bucketNames []string) *minio.Client {
	minioClient, err := ConnectMinio(cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to minio")
	}

	for _, bucketName := range bucketNames {
		err = CreateBucketIfNotExists(context.Background(), minioClient, bucketName)
		if err != nil {
			l.Fatal().Str("bucketName", bucketName).Err(err).Msg("Failed to create bucket")
		}
	}

	l.Trace().Msg("Successfully connected to MinIO")
	return minioClient
}

func ConnectMinio(cfg *ConnectConfig) (*minio.Client, error) {
	// Подключение к minio
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Пинг
	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		return nil, fmt.Errorf("MinIO ping failed: %w", err)
	}

	return minioClient, nil
}

func CreateBucketIfNotExists(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	// Создаем bucket
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Проверяем существует ли bucket
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return fmt.Errorf("create bucket error: %s", err)
		}
	}

	return nil
}
