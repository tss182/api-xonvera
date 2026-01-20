package services

import (
	"context"
	"testing"
	"time"

	"app/xonvera-core/internal/core/domain"
	"app/xonvera-core/internal/infrastructure/config"
	"app/xonvera-core/internal/testutil"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		ExistsByPhoneFunc: func(ctx context.Context, phone string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, user *domain.User) error {
			user.ID = 1
			return nil
		},
	}

	mockTokenRepo := &testutil.MockTokenRepository{
		CreateFunc: func(ctx context.Context, token *domain.Token) error {
			return nil
		},
	}

	mockTokenService := &testutil.MockTokenService{
		GenerateTokenPairFunc: func(userID uint) (*domain.TokenPair, error) {
			return &domain.TokenPair{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				ExpiresAt:    time.Now().Unix(),
			}, nil
		},
	}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	req := &domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	response, err := service.Register(context.Background(), req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if response == nil {
		t.Error("expected response, got nil")
	}
	if response.User.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", response.User.Email)
	}
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
	}

	mockTokenRepo := &testutil.MockTokenRepository{}
	mockTokenService := &testutil.MockTokenService{}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	req := &domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	response, err := service.Register(context.Background(), req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if response != nil {
		t.Errorf("expected nil response, got %v", response)
	}
}

func TestAuthService_Register_PhoneAlreadyExists(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		ExistsByPhoneFunc: func(ctx context.Context, phone string) (bool, error) {
			return true, nil
		},
	}

	mockTokenRepo := &testutil.MockTokenRepository{}
	mockTokenService := &testutil.MockTokenService{}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	req := &domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	response, err := service.Register(context.Background(), req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if response != nil {
		t.Errorf("expected nil response, got %v", response)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: string(hashedPassword),
	}

	mockUserRepo := &testutil.MockUserRepository{
		FindByEmailOrPhoneFunc: func(ctx context.Context, username string) (*domain.User, error) {
			return existingUser, nil
		},
	}

	mockTokenRepo := &testutil.MockTokenRepository{
		CreateFunc: func(ctx context.Context, token *domain.Token) error {
			return nil
		},
	}

	mockTokenService := &testutil.MockTokenService{
		GenerateTokenPairFunc: func(userID uint) (*domain.TokenPair, error) {
			return &domain.TokenPair{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				ExpiresAt:    time.Now().Unix(),
			}, nil
		},
	}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	req := &domain.LoginRequest{
		Username: "john@example.com",
		Password: "password123",
	}

	response, err := service.Login(context.Background(), req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if response == nil {
		t.Error("expected response, got nil")
	}
	if response.User.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", response.User.Email)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: string(hashedPassword),
	}

	mockUserRepo := &testutil.MockUserRepository{
		FindByEmailOrPhoneFunc: func(ctx context.Context, username string) (*domain.User, error) {
			return existingUser, nil
		},
	}

	mockTokenRepo := &testutil.MockTokenRepository{}
	mockTokenService := &testutil.MockTokenService{}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	req := &domain.LoginRequest{
		Username: "john@example.com",
		Password: "wrongpassword",
	}

	response, err := service.Login(context.Background(), req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if response != nil {
		t.Errorf("expected nil response, got %v", response)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{
		FindByEmailOrPhoneFunc: func(ctx context.Context, username string) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	mockTokenRepo := &testutil.MockTokenRepository{}
	mockTokenService := &testutil.MockTokenService{}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	req := &domain.LoginRequest{
		Username: "nonexistent@example.com",
		Password: "password123",
	}

	response, err := service.Login(context.Background(), req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if response != nil {
		t.Errorf("expected nil response, got %v", response)
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{}

	mockTokenRepo := &testutil.MockTokenRepository{
		FindByAccessTokenFunc: func(ctx context.Context, accessToken string) (*domain.Token, error) {
			return &domain.Token{
				ID:     1,
				UserID: 1,
			}, nil
		},
		DeleteFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	mockTokenService := &testutil.MockTokenService{}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	err := service.Logout(context.Background(), "valid_token")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAuthService_Logout_InvalidToken(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{}

	mockTokenRepo := &testutil.MockTokenRepository{
		FindByAccessTokenFunc: func(ctx context.Context, accessToken string) (*domain.Token, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	mockTokenService := &testutil.MockTokenService{}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	err := service.Logout(context.Background(), "invalid_token")

	if err != nil {
		t.Errorf("expected no error (Logout returns nil for non-existent tokens), got %v", err)
	}
}

func TestAuthService_ValidateAccessToken_Success(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepository{}

	mockTokenRepo := &testutil.MockTokenRepository{
		FindByAccessTokenFunc: func(ctx context.Context, accessToken string) (*domain.Token, error) {
			return &domain.Token{
				ID:     1,
				UserID: 1,
			}, nil
		},
	}

	mockTokenService := &testutil.MockTokenService{
		ValidateTokenFunc: func(token string) (uint, error) {
			return 1, nil
		},
	}

	tokenConfig := &config.TokenConfig{
		SecretKey:         "test_secret_key",
		Expired:       time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewAuthService(mockUserRepo, mockTokenRepo, mockTokenService, tokenConfig)

	userID, err := service.ValidateAccessToken(context.Background(), "valid_token")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if userID != 1 {
		t.Errorf("expected user ID 1, got %d", userID)
	}
}
