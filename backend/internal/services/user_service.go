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
	if bcrypt.CompareHashAndPassword([]byte(userWithLogin.PasswordHash), []byte(userLogin.Password)) == nil {
		return userWithLogin.ID, nil
	}
	return 0, errors.New("invalid password")
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*entities.UserGet, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*entities.UserGet, error) {
	users, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUserLogin(ctx context.Context, userUpdateLoginRequest *entities.UserUpdateLoginRequest) error {
	err := s.userRepository.UpdateUserLogin(ctx, userUpdateLoginRequest.ID, userUpdateLoginRequest.Login)
	return err
}

func (s *UserService) UpdateUserPassword(ctx context.Context, userUpdatePasswordRequest *entities.UserUpdatePasswordRequest) error {
	// Хешируем новый пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userUpdatePasswordRequest.Password), bcryptCost)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdateUserPassword(ctx, userUpdatePasswordRequest.ID, string(passwordHash))
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepository.DeleteUser(ctx, id)
	return err
}
