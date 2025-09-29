package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/repositories"
	"time"
)

const (
	AccessTokenType  = "access_token"
	RefreshTokenType = "refresh_token"
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

func (s *AuthService) RegistrationUser(ctx context.Context, userCreate *entities.UserCreateRequest) (*entities.TokenPair, error) {
	// Создаем пользователя
	userID, err := s.userService.CreateUser(ctx, userCreate)
	if err != nil {
		return nil, err
	}

	// Генерируем пару access и refresh токенов
	tokenPair, err := s.CreateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh токен в кеш (если не получилось - не критично)
	_ = s.tokenRepository.AddRefreshToken(ctx, userID, tokenPair.RefreshToken, s.refreshTokenLifetime)

	return tokenPair, nil
}

func (s *AuthService) LoginUser(ctx context.Context, userLogin *entities.UserLoginRequest) (*entities.TokenPair, error) {
	// Проверяем, есть ли пользователь с таким логином в системе
	userID, err := s.userService.LoginUser(ctx, userLogin)
	if err != nil {
		return nil, err
	}

	// Генерируем токены
	tokenPair, err := s.CreateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh токен в кеш (если не получилось - не критично)
	_ = s.tokenRepository.AddRefreshToken(ctx, userID, tokenPair.RefreshToken, s.refreshTokenLifetime)

	return tokenPair, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error) {
	// Получаем данные из refresh токена
	refreshTokenClaims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("incorrect refresh token")
	}

	// Проверка типа токена
	if refreshTokenClaims.Type != RefreshTokenType {
		return nil, errors.New("not a refresh token")
	}

	// Проверяем, существует ли этот токен
	err = s.tokenRepository.CheckExistsRefreshToken(ctx, refreshTokenClaims.UserID, refreshToken)
	if err != nil {
		return nil, err
	}

	// Создаем новый access токен
	newAccessToken, err := s.GenerateToken(refreshTokenClaims.UserID, true)
	if err != nil {
		return nil, errors.New("failed to create access token")
	}

	// Создаем новый refresh токен
	newRefreshToken, err := s.GenerateToken(refreshTokenClaims.UserID, false)
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	// Заменяем старый refresh токен на новый (если не получилось - не критично)
	_ = s.tokenRepository.ReplaceRefreshToken(ctx, refreshTokenClaims.UserID, refreshToken, newRefreshToken, s.refreshTokenLifetime)

	return &entities.TokenPair{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}

// Генерация пары access и refresh токенов
func (s *AuthService) CreateTokens(userID int) (*entities.TokenPair, error) {
	// Генерируем access токен
	accessToken, err := s.GenerateToken(userID, true)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh токен
	refreshToken, err := s.GenerateToken(userID, false)
	if err != nil {
		return nil, err
	}

	// Возвращаем оба токена
	return &entities.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Создание jwt токена
func (s *AuthService) GenerateToken(userID int, isAccessToken bool) (string, error) {
	// Время создания токена
	now := time.Now()

	// Определяем тип токена
	tokenLifetime := s.accessTokenLifetime
	tokenType := AccessTokenType
	if !isAccessToken {
		tokenLifetime = s.refreshTokenLifetime
		tokenType = RefreshTokenType
	}

	// Тело токена
	claims := &entities.TokenClaims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
			ID:        fmt.Sprintf("token-%d", s.tokenID),
		},
	}
	s.tokenID++

	// Подписываем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Парсинг jwt токена
func (s *AuthService) ParseToken(tokenString string) (*entities.TokenClaims, error) {
	// Подтверждение подлинности токена
	token, err := jwt.ParseWithClaims(tokenString, &entities.TokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, errors.New("can't verify token")
	}

	// Парсинг тела токена
	if claims, ok := token.Claims.(*entities.TokenClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
