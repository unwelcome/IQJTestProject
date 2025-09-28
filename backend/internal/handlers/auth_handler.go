package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unwelcome/iqjtest/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return nil
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return nil
}
