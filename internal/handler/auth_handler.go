package handler

import (
	"trainers-backend/internal/dto"
	"trainers-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) TelegramAuth(c *fiber.Ctx) error {
	var req dto.TelegramAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:  "error",
			Message: "invalid request body",
		})
	}

	user, token, err := h.authService.AuthenticateTelegram(c.Context(), req.InitData)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.Response{
		Status: "success",
		Data: dto.AuthResponse{
			Token: token,
			User:  *user,
		},
	})
}
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	// В реальном проекте здесь был бы запрос к userRepo
	// Для краткости вернем то, что есть в токене
	return c.JSON(dto.Response{
		Status: "success",
		Data: fiber.Map{
			"user_id": userID,
			"role":    c.Locals("user_role"),
		},
	})
}

// Dashboard Handler
func (h *AuthHandler) GetDashboard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("user_role").(string)

	// Заглушка для дашборда. Позже добавим реальную статистику из БД
	return c.JSON(dto.Response{
		Status: "success",
		Data: fiber.Map{
			"user_id": userID,
			"role":    role,
			"stats":   fiber.Map{"clients": 0, "workouts": 0},
		},
	})
}
