package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/repositories"
	"strconv"
	"time"
)

type AuthService struct {
	tokenRepository      *repositories.AuthRepository
	secretKey            string
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
	tokenID              int
}

func NewAuthService(tokenRepository *repositories.AuthRepository, secretKey string, accessTokenLifetime time.Duration, refreshTokenLifetime time.Duration) *AuthService {
	return &AuthService{
		tokenRepository:      tokenRepository,
		secretKey:            secretKey,
		accessTokenLifetime:  accessTokenLifetime,
		refreshTokenLifetime: refreshTokenLifetime,
		tokenID:              1,
	}
}

func (s *AuthService) CreateTokens(ctx context.Context, userID int) (*entities.TokenPair, error) {
	// Генерируем access токен
	accessToken, err := s.CreateToken(ctx, strconv.Itoa(userID), s.accessTokenLifetime)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh токен
	refreshToken, err := s.CreateToken(ctx, strconv.Itoa(userID), s.refreshTokenLifetime)
	if err != nil {
		return nil, err
	}

	// Возвращаем оба токена
	return &entities.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error) {
	// Получаем данные из refresh токена
	refreshTokenClaims, err := s.ParseToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Создаем новый access токен
	newAccessToken, err := s.CreateToken(ctx, refreshTokenClaims.Subject, s.accessTokenLifetime)
	if err != nil {
		return nil, err
	}

	// Создаем новый refresh токен
	newRefreshToken, err := s.CreateToken(ctx, refreshTokenClaims.Subject, s.refreshTokenLifetime)
	if err != nil {
		return nil, err
	}

	return &entities.TokenPair{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}

// Создание jwt токена
func (s *AuthService) CreateToken(ctx context.Context, userID string, tokenLifetime time.Duration) (string, error) {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer:    "iqjtest-auth-server",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"iqjtest-api"},
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        fmt.Sprintf("token-%d", s.tokenID),
	}
	s.tokenID++

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Парсинг jwt токена
func (s *AuthService) ParseToken(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error) {
	parsedToken, _ := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
