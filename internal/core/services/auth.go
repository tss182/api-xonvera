package services

import (
	"context"
	"errors"
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
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("failed to check email existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("400:email already registered")
	}

	// Check if phone already exists
	exists, err = s.userRepo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		logger.Error("failed to check phone existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("400:phone number already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", zap.Error(err))
		return nil, errors.New("400:invalid register request")
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return nil, errors.New("400:invalid register request")
	}

	// Generate token pair
	tokenPair, err := s.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		logger.Error("failed to generate token pair", zap.Error(err))
		return nil, err
	}

	// Save token to database
	if err := s.saveToken(ctx, user.ID, tokenPair); err != nil {
		logger.Error("failed to save token", zap.Error(err))
		return nil, err
	}

	logger.Info("user registered successfully", zap.Uint("user_id", user.ID), zap.String("email", user.Email))

	return &domain.AuthResponse{
		User:         user,
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
			return nil, errors.New("400:invalid credentials")
		}
		logger.Error("failed to find user", zap.Error(err))
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("400:invalid credentials")
	}

	// Generate token pair
	tokenPair, err := s.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		logger.Error("failed to generate token pair", zap.Error(err))
		return nil, err
	}

	// Save token to database
	if err := s.saveToken(ctx, user.ID, tokenPair); err != nil {
		logger.Error("failed to save token", zap.Error(err))
		return nil, err
	}

	logger.Info("user logged in successfully", zap.Uint("user_id", user.ID))

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	// Validate refresh token format
	userID, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("400:invalid refresh token")
	}

	// Find token in database
	storedToken, err := s.tokenRepo.FindByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, errors.New("token not found")) {
			return nil, errors.New("400:invalid refresh token")
		}
		logger.Error("failed to find refresh token", zap.Error(err))
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(storedToken.RefreshExpiresAt) {
		return nil, errors.New("400:refresh token has expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Error("failed to find user", zap.Error(err))
		return nil, errors.New("400:invalid refresh token")
	}

	// Generate new token pair
	tokenPair, err := s.tokenService.GenerateTokenPair(userID)
	if err != nil {
		logger.Error("failed to generate token pair", zap.Error(err))
		return nil, err
	}

	// Update token in database
	storedToken.AccessToken = tokenPair.AccessToken
	storedToken.RefreshToken = tokenPair.RefreshToken
	storedToken.ExpiresAt = time.Unix(tokenPair.ExpiresAt, 0)
	storedToken.RefreshExpiresAt = time.Now().Add(s.tokenConfig.RefreshExpired)

	if err := s.tokenRepo.Update(ctx, storedToken); err != nil {
		logger.Error("failed to update token", zap.Error(err))
		return nil, err
	}

	logger.Info("token refreshed successfully", zap.Uint("user_id", userID))

	return &domain.AuthResponse{
		User:         user,
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
		logger.Debug("token not found during logout, already expired or deleted")
		return nil
	}

	// Delete token from repository (will work for both DB and Redis implementations)
	if err := s.tokenRepo.DeleteByUserID(ctx, storedToken.UserID); err != nil {
		logger.Error("failed to delete token", zap.Error(err))
		return err
	}

	logger.Info("user logged out successfully", zap.Uint("user_id", storedToken.UserID))
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
		return 0, err
	}

	// Then check if token exists in database (not logged out)
	_, err = s.tokenRepo.FindByAccessToken(ctx, accessToken)
	if err != nil {
		return 0, errors.New("token has been invalidated")
	}

	return userID, nil
}
