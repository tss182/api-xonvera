package http

import (
	"context"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/utils/validator"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service portService.AuthService
	rto     time.Duration
}

func NewAuthHandler(service portService.AuthService, rto time.Duration) *AuthHandler {
	return &AuthHandler{
		service: service,
		rto:     rto,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with name, email, phone, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Router /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	var req dto.RegisterRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in register", zap.Strings("error validation query register", err))
		return BadRequest(c, err)
	}

	res, err := h.service.Register(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, res)
}

// Login handles user login
// @Summary Login user
// @Description Login with email or phone and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Failure 401 {object} Resp
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	var req dto.LoginRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in login", zap.Strings("error validation query login", err))
		return BadRequest(c, err)
	}

	res, err := h.service.Login(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, res)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Failure 401 {object} Resp
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	var req dto.RefreshTokenRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in refresh token", zap.Strings("error validation query refresh token", err))
		return BadRequest(c, err)
	}

	res, err := h.service.RefreshToken(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, res)
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate current access token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Resp
// @Failure 401 {object} Resp
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	// Get access token from context (set by auth middleware)
	accessToken, ok := c.Locals("accessToken").(string)
	if !ok || accessToken == "" {
		return NoAuth(c)
	}

	err := h.service.Logout(ctx, accessToken)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, nil)
}

// fiber:context-methods migrated
