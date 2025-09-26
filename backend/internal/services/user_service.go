package services

import (
	"context"
	"errors"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

const bcryptCost = 10

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, userCreate *entities.UserCreateRequest) (*entities.UserCreateResponse, error) {
	// Переводим пароль из строки в срез байт
	bytePassword := []byte(userCreate.Password)
	if len(bytePassword) >= 70 { // Больше 72 байт библиотека не захеширует
		return nil, errors.New("password too long")
	}

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcryptCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &entities.User{Login: userCreate.Login, PasswordHash: string(passwordHash)}

	// Добавляем пользователя в бд
	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &entities.UserCreateResponse{ID: user.ID}, nil
}
