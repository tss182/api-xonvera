package http

import (
	"context"
	"time"

	"app/xonvera-core/internal/core/domain"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/utils/validator"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService    portService.AuthService
	requestTimeout time.Duration
}

func NewAuthHandler(authService portService.AuthService, requestTimeout time.Duration) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		requestTimeout: requestTimeout,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with name, email, phone, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register Request"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.requestTimeout)
	defer cancel()

	var req domain.RegisterRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in auth service", zap.Strings("error validation query", err))
		return BadRequest(c, err)
	}

	response, err := h.authService.Register(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, response)
}

// Login handles user login
// @Summary Login user
// @Description Login with email or phone and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login Request"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.requestTimeout)
	defer cancel()

	var req domain.LoginRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in auth service", zap.Strings("error validation query", err))
		return BadRequest(c, err)
	}

	response, err := h.authService.Login(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, response)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.requestTimeout)
	defer cancel()

	var req domain.RefreshTokenRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in auth service", zap.Strings("error validation query", err))
		return BadRequest(c, err)
	}

	response, err := h.authService.RefreshToken(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, response)
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate current access token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.requestTimeout)
	defer cancel()

	// Get access token from context (set by auth middleware)
	accessToken, ok := c.Locals("accessToken").(string)
	if !ok || accessToken == "" {
		return NoAuth(c)
	}

	err := h.authService.Logout(ctx, accessToken)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, nil)
}
