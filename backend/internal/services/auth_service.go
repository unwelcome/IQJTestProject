package services

import (
	"context"
	"fmt"
	"github.com/unwelcome/iqjtest/pkg/utils"
	"time"

	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/repositories"
)

type AuthService struct {
	userService     *UserService
	tokenRepository *repositories.AuthRepository

	secretKey            string
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration

	tokenID int
}

func NewAuthService(userService *UserService, tokenRepository *repositories.AuthRepository, secretKey string, accessTokenLifetime time.Duration, refreshTokenLifetime time.Duration) *AuthService {
	return &AuthService{
		userService:     userService,
		tokenRepository: tokenRepository,

		secretKey:            secretKey,
		accessTokenLifetime:  accessTokenLifetime,
		refreshTokenLifetime: refreshTokenLifetime,

		tokenID: 1,
	}
}

func (s *AuthService) RegistrationUser(ctx context.Context, userCreate *entities.UserCreateRequest) (*entities.AuthResponse, error) {

	// Создаем пользователя
	userID, err := s.userService.CreateUser(ctx, userCreate)
	if err != nil {
		return nil, err
	}

	// Генерируем пару access и refresh токенов
	tokenPair, err := utils.CreateTokens(userID, s.secretKey, s.accessTokenLifetime, s.refreshTokenLifetime, &s.tokenID)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}

	// Сохраняем refresh токен в кеш (если не получилось - не критично)
	_ = s.tokenRepository.AddToken(ctx, userID, tokenPair.RefreshToken, entities.RefreshTokenType, s.refreshTokenLifetime)

	return &entities.AuthResponse{TokenPair: tokenPair, UserID: userID}, nil
}

func (s *AuthService) LoginUser(ctx context.Context, userLogin *entities.UserLoginRequest) (*entities.AuthResponse, error) {

	// Проверяем, есть ли пользователь с таким логином в системе и получаем его ID
	userID, err := s.userService.LoginUser(ctx, userLogin)
	if err != nil {
		return nil, err
	}

	// Генерируем токены
	tokenPair, err := utils.CreateTokens(userID, s.secretKey, s.accessTokenLifetime, s.refreshTokenLifetime, &s.tokenID)
	if err != nil {
		return nil, fmt.Errorf("login user error: %w", err)
	}

	// Сохраняем refresh токен в кеш (если не получилось - не критично)
	_ = s.tokenRepository.AddToken(ctx, userID, tokenPair.RefreshToken, entities.RefreshTokenType, s.refreshTokenLifetime)

	return &entities.AuthResponse{TokenPair: tokenPair, UserID: userID}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error) {

	// Парсим refresh токен
	tokenClaims, err := utils.ParseToken(refreshToken, s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("refresh tokens error: %w", err)
	}

	// Проверяем тип токена
	if tokenClaims.Type != entities.RefreshTokenType {
		return nil, fmt.Errorf("refresh tokens error: invalid token type")
	}

	// Создаем новую пару токенов
	tokenPair, err := utils.CreateTokens(tokenClaims.UserID, s.secretKey, s.accessTokenLifetime, s.refreshTokenLifetime, &s.tokenID)
	if err != nil {
		return nil, fmt.Errorf("refresh tokens error: %w", err)
	}

	// Заменяем старый refresh токен на новый
	err = s.tokenRepository.ReplaceToken(ctx, tokenClaims.UserID, refreshToken, tokenPair.RefreshToken, entities.RefreshTokenType, s.refreshTokenLifetime)
	if err != nil {
		return nil, fmt.Errorf("refresh tokens error: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) DeleteRefreshToken(ctx context.Context, userID int, refreshToken string) error {

	// Удаляем токен
	err := s.tokenRepository.DeleteToken(ctx, userID, refreshToken, entities.RefreshTokenType)
	if err != nil {
		return fmt.Errorf("delete refresh token error: %w", err)
	}

	return nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userID int) error {

	// Удаляем пользователя из бд
	err := s.userService.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	// Удаляем все токены пользователя
	err = s.tokenRepository.DeleteAllTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete all user tokens error: %w", err)
	}

	return nil
}
