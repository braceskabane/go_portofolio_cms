package handler

import (
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/middleware"
	"portfolio-cms/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return BadRequest(c, err.Error())
	}
	return Created(c, user, "User registered successfully")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}
	user, err := h.authService.GetByEmail(req.Email, req.Password)
	if err != nil {
		return Unauthorized(c, err.Error())
	}
	// Generate real tokens here in handler (avoids import cycle)
	accessToken, err := middleware.GenerateAccessToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return InternalError(c, "Failed to generate token")
	}
	refreshToken, err := middleware.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return InternalError(c, "Failed to generate token")
	}
	return OK(c, dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "Login successful")
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return BadRequest(c, err.Error())
	}

	// Parse refresh token to get userID
	userID, err := middleware.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return Unauthorized(c, "Invalid or expired refresh token")
	}

	// Get user to verify still exists and get updated role
	user, err := h.authService.GetByID(userID)
	if err != nil {
		return NotFound(c, "User not found")
	}

	// Generate new tokens
	accessToken, err := middleware.GenerateAccessToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return InternalError(c, "Failed to generate token")
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return InternalError(c, "Failed to generate token")
	}

	return OK(c, dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "Token refreshed successfully")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	user, err := h.authService.GetByID(userID)
	if err != nil {
		return NotFound(c, "User not found")
	}
	return OK(c, user)
}
