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

func (s *UserService) CreateUser(ctx context.Context, userCreate *entities.UserCreateRequest) (int, error) {
	// Переводим пароль из строки в срез байт
	bytePassword := []byte(userCreate.Password)
	if len(bytePassword) >= 70 { // Больше 72 байт библиотека не захеширует
		return 0, errors.New("password too long")
	}

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcryptCost)
	if err != nil {
		return 0, err
	}

	// Создаем пользователя
	user := &entities.User{Login: userCreate.Login, PasswordHash: string(passwordHash)}

	// Добавляем пользователя в бд
	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *UserService) LoginUser(ctx context.Context, userLogin *entities.UserLoginRequest) (int, error) {
	// Получаем пользователя с данным логином
	userWithLogin, err := s.userRepository.GetUserByLogin(ctx, userLogin.Login)
	if err != nil {
		return 0, errors.New("user not found")
	}

	// Проверяем пароль
	if bcrypt.CompareHashAndPassword([]byte(userWithLogin.PasswordHash), []byte(userLogin.Password)) != nil {
		return 0, errors.New("invalid password")
	}

	return userWithLogin.ID, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*entities.UserGet, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.ID = userID
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*entities.UserGet, error) {
	users, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, userID int, userUpdatePasswordRequest *entities.UserUpdatePasswordRequest) error {
	// Хешируем новый пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userUpdatePasswordRequest.Password), bcryptCost)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdateUserPassword(ctx, userID, string(passwordHash))
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	err := s.userRepository.DeleteUser(ctx, userID)
	return err
}
