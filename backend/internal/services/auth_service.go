package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *AuthService) RegistrationUser(ctx context.Context, userCreate *entities.UserCreateRequest) (*entities.TokenPair, error) {
	// Создаем пользователя
	userID, err := s.userService.CreateUser(ctx, userCreate)
	if err != nil {
		return nil, err
	}

	// Генерируем пару access и refresh токенов
	tokenPair, err := CreateTokens(userID, s)
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh токен в кеш (если не получилось - не критично)
	_ = s.tokenRepository.AddToken(ctx, userID, tokenPair.RefreshToken, entities.RefreshTokenType, s.refreshTokenLifetime)

	return tokenPair, nil
}

func (s *AuthService) LoginUser(ctx context.Context, userLogin *entities.UserLoginRequest) (*entities.TokenPair, error) {
	// Проверяем, есть ли пользователь с таким логином в системе
	userID, err := s.userService.LoginUser(ctx, userLogin)
	if err != nil {
		return nil, err
	}

	// Генерируем токены
	tokenPair, err := CreateTokens(userID, s)
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh токен в кеш (если не получилось - не критично)
	_ = s.tokenRepository.AddToken(ctx, userID, tokenPair.RefreshToken, entities.RefreshTokenType, s.refreshTokenLifetime)

	return tokenPair, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error) {
	// Парсим refresh токен
	tokenClaims, err := ParseToken(refreshToken, s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Проверяем тип токена
	if tokenClaims.Type != entities.RefreshTokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	// Создаем новую пару токенов
	tokenPair, err := CreateTokens(tokenClaims.UserID, s)
	if err != nil {
		return nil, errors.New("failed to create new tokens")
	}

	// Заменяем старый refresh токен на новый
	err = s.tokenRepository.ReplaceToken(ctx, tokenClaims.UserID, refreshToken, tokenPair.RefreshToken, entities.RefreshTokenType, s.refreshTokenLifetime)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (int, error) {
	tokenClaims, err := ParseToken(accessToken, s.secretKey)
	if err != nil {
		return 0, err
	}

	if tokenClaims.Type != entities.AccessTokenType {
		return 0, errors.New("invalid token type")
	}

	return tokenClaims.UserID, nil
}

func (s *AuthService) DeleteRefreshToken(ctx context.Context, userID int, refreshToken string) error {
	err := s.tokenRepository.DeleteToken(ctx, userID, refreshToken, entities.RefreshTokenType)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// Генерация пары access и refresh токенов
func CreateTokens(userID int, s *AuthService) (*entities.TokenPair, error) {
	// Генерируем access токен
	accessToken, err := GenerateToken(userID, s.secretKey, entities.AccessTokenType, s.accessTokenLifetime, s.tokenID)
	s.tokenID++
	if err != nil {
		return nil, err
	}

	// Генерируем refresh токен
	refreshToken, err := GenerateToken(userID, s.secretKey, entities.RefreshTokenType, s.refreshTokenLifetime, s.tokenID)
	s.tokenID++
	if err != nil {
		return nil, err
	}

	// Возвращаем оба токена
	return &entities.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Создание jwt токена
func GenerateToken(userID int, secretKey string, tokenType string, tokenLifetime time.Duration, tokenID int) (string, error) {
	// Время создания токена
	now := time.Now()

	// Тело токена
	claims := &entities.TokenClaims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
			ID:        fmt.Sprintf("token-%d", tokenID),
		},
	}

	// Подписываем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Парсинг jwt токена
func ParseToken(tokenString string, secretKey string) (*entities.TokenClaims, error) {
	// Подтверждение подлинности токена
	token, err := jwt.ParseWithClaims(tokenString, &entities.TokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("can't verify token")
	}

	// Парсинг тела токена
	if claims, ok := token.Claims.(*entities.TokenClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
