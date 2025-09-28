package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/entities"
	"github.com/unwelcome/iqjtest/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUserByID
// @Summary получение пользователя по ID
// @Description получаем пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.UserGet
// @Failure 400 {object} entities.ErrorEntity
// @Failure 404 {object} entities.ErrorEntity
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// Ограничение времени выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Missing id"})
	}

	if userID < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(200).JSON(user)
}

// GetAllUsers
// @Summary получение всех пользователей
// @Description получение всех пользователей
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} []entities.UserGet
// @Failure 500 {object} entities.ErrorEntity
// @Router /user/all [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(users)
}

// UpdateUserLogin
// @Summary обновление логина пользователя
// @Description обновление логина пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.UserUpdateLoginRequest true "Данные пользователя"
// @Success 200 {object} entities.UserUpdateLoginResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /user/login [patch]
func (h *UserHandler) UpdateUserLogin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userUpdateLoginRequest := &entities.UserUpdateLoginRequest{}
	if err := c.BodyParser(&userUpdateLoginRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := h.userService.UpdateUserLogin(ctx, userUpdateLoginRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(&entities.UserUpdateLoginResponse{ID: userUpdateLoginRequest.ID})
}

// UpdateUserPassword
// @Summary обновление пароля пользователя
// @Description обновление пароля пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.UserUpdatePasswordRequest true "Данные пользователя"
// @Success 200 {object} entities.UserUpdatePasswordResponse
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /user/password [patch]
func (h *UserHandler) UpdateUserPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userUpdatePasswordRequest := &entities.UserUpdatePasswordRequest{}
	if err := c.BodyParser(&userUpdatePasswordRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := h.userService.UpdateUserPassword(ctx, userUpdatePasswordRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(&entities.UserUpdatePasswordResponse{ID: userUpdatePasswordRequest.ID})
}

//TODO
// убрать DeleteUser в AuthHandler

// DeleteUser
// @Summary удаление пользователя
// @Description удаление пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Failure 400 {object} entities.ErrorEntity
// @Failure 500 {object} entities.ErrorEntity
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Missing id"})
	}

	if userID < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	err = h.userService.DeleteUser(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).SendString("User deleted")
}
