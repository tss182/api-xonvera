package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/config"
	"app/xonvera-core/internal/infrastructure/logger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	userRepo     portRepository.UserRepository
	tokenRepo    portRepository.TokenRepository
	tokenService portService.TokenService
	tokenConfig  *config.TokenConfig
}

func NewAuthService(
	userRepo portRepository.UserRepository,
	tokenRepo portRepository.TokenRepository,
	tokenService portService.TokenService,
	tokenConfig *config.TokenConfig,
) portService.AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
		tokenConfig:  tokenConfig,
	}
}

func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if phone already exists
	exists, err := s.userRepo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		logger.StdContextError(ctx, "failed to check phone existence", zap.Error(err), zap.String("phone", req.Phone))
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("400:phone number already registered")
	}

	// Hash password with secure cost factor
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.StdContextError(ctx, "failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("400:invalid register request")
	}

	// Create new user
	user := &domain.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.StdContextError(ctx, "failed to create user", zap.Error(err), zap.String("phone", req.Phone))
		return nil, fmt.Errorf("400:invalid register request")
	}

	// Generate token pair
	tokenPair, err := s.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		logger.Error("failed to generate token pair", zap.Error(err), zap.Uint("user_id", user.ID))
		return nil, err
	}

	// Save token to database
	if err := s.saveToken(ctx, user.ID, tokenPair); err != nil {
		logger.Error("failed to save token", zap.Error(err), zap.Uint("user_id", user.ID))
		return nil, err
	}

	logger.StdContextInfo(ctx, "user registered successfully",
		zap.Uint("user_id", user.ID),
		zap.String("phone", user.Phone),
	)

	return &domain.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email or phone
	user, err := s.userRepo.FindByEmailOrPhone(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.StdContextDebug(ctx, "login failed: invalid credentials", zap.String("username", req.Username))
			return nil, fmt.Errorf("400:invalid credentials")
		}
		logger.StdContextError(ctx, "failed to find user", zap.Error(err), zap.String("username", req.Username))
		return nil, err
	}

	// Verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.StdContextDebug(ctx, "login failed: password mismatch", zap.Uint("user_id", user.ID))
		return nil, fmt.Errorf("400:invalid credentials")
	}

	// Generate token pair
	tokenPair, err := s.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		logger.StdContextError(ctx, "failed to generate token pair", zap.Error(err), zap.Uint("user_id", user.ID))
		return nil, err
	}

	// Save token to database
	if err := s.saveToken(ctx, user.ID, tokenPair); err != nil {
		logger.StdContextError(ctx, "failed to save token", zap.Error(err), zap.Uint("user_id", user.ID))
		return nil, err
	}

	logger.StdContextInfo(ctx, "user logged in successfully",
		zap.Uint("user_id", user.ID),
		zap.String("phone", user.Phone),
	)

	return &domain.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	// Validate refresh token format
	userID, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("400:invalid refresh token")
	}

	// Find token in database
	storedToken, err := s.tokenRepo.FindByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, fmt.Errorf("404:not found token")) {
			return nil, fmt.Errorf("400:invalid refresh token")
		}
		logger.StdContextError(ctx, "failed to find refresh token", zap.Error(err))
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(storedToken.RefreshExpiresAt) {
		return nil, fmt.Errorf("400:refresh token has expired")
	}

	// Generate new token pair
	tokenPair, err := s.tokenService.GenerateTokenPair(userID)
	if err != nil {
		logger.StdContextError(ctx, "failed to generate token pair", zap.Error(err))
		return nil, err
	}

	// Update token in database
	storedToken.AccessToken = tokenPair.AccessToken
	storedToken.RefreshToken = tokenPair.RefreshToken
	storedToken.ExpiresAt = time.Unix(tokenPair.ExpiresAt, 0)
	storedToken.RefreshExpiresAt = time.Now().Add(s.tokenConfig.RefreshExpired)

	if err := s.tokenRepo.Update(ctx, storedToken); err != nil {
		logger.StdContextError(ctx, "failed to update token", zap.Error(err))
		return nil, err
	}

	logger.StdContextInfo(ctx, "token refreshed successfully", zap.Uint("user_id", userID))

	return &domain.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string) error {
	// Find token by access token to verify it exists
	storedToken, err := s.tokenRepo.FindByAccessToken(ctx, accessToken)
	if err != nil {
		// Token might already be deleted, still return success
		logger.StdContextDebug(ctx, "token not found during logout, already expired or deleted")
		return nil
	}

	// Delete token from repository (will work for both DB and Redis implementations)
	if err := s.tokenRepo.DeleteByUserID(ctx, storedToken.UserID); err != nil {
		logger.StdContextError(ctx, "failed to delete token", zap.Error(err))
		return err
	}

	logger.StdContextInfo(ctx, "user logged out successfully", zap.Uint("user_id", storedToken.UserID))
	return nil
}

func (s *authService) saveToken(ctx context.Context, userID uint, tokenPair *domain.TokenPair) error {
	token := &domain.Token{
		UserID:           userID,
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresAt:        time.Unix(tokenPair.ExpiresAt, 0),
		RefreshExpiresAt: time.Now().Add(s.tokenConfig.RefreshExpired),
	}

	return s.tokenRepo.Create(ctx, token)
}

func (s *authService) ValidateAccessToken(ctx context.Context, accessToken string) (uint, error) {
	// First validate token format and expiration
	userID, err := s.tokenService.ValidateToken(accessToken)
	if err != nil {
		logger.StdContextDebug(ctx, "token validation failed", zap.Error(err))
		return 0, err
	}

	// Then check if token exists in database (not logged out)
	_, err = s.tokenRepo.FindByAccessToken(ctx, accessToken)
	if err != nil {
		logger.StdContextDebug(ctx, "token not found in repository (possibly logged out)", zap.Uint("user_id", userID))
		return 0, fmt.Errorf("token has been invalidated")
	}

	return userID, nil
}
