package testutil

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

// MockUserRepository is a mock implementation of ports.UserRepository for testing
type MockUserRepository struct {
	CreateFunc             func(ctx context.Context, user *domain.User) error
	FindByEmailFunc        func(ctx context.Context, email string) (*domain.User, error)
	FindByPhoneFunc        func(ctx context.Context, phone string) (*domain.User, error)
	FindByEmailOrPhoneFunc func(ctx context.Context, username string) (*domain.User, error)
	FindByIDFunc           func(ctx context.Context, id uint) (*domain.User, error)
	ExistsByEmailFunc      func(ctx context.Context, email string) (bool, error)
	ExistsByPhoneFunc      func(ctx context.Context, phone string) (bool, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	if m.FindByPhoneFunc != nil {
		return m.FindByPhoneFunc(ctx, phone)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByEmailOrPhone(ctx context.Context, username string) (*domain.User, error) {
	if m.FindByEmailOrPhoneFunc != nil {
		return m.FindByEmailOrPhoneFunc(ctx, username)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(ctx, email)
	}
	return false, nil
}

func (m *MockUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	if m.ExistsByPhoneFunc != nil {
		return m.ExistsByPhoneFunc(ctx, phone)
	}
	return false, nil
}

// MockTokenRepository is a mock implementation of ports.TokenRepository for testing
type MockTokenRepository struct {
	CreateFunc              func(ctx context.Context, token *domain.Token) error
	FindByAccessTokenFunc   func(ctx context.Context, accessToken string) (*domain.Token, error)
	FindByRefreshTokenFunc  func(ctx context.Context, refreshToken string) (*domain.Token, error)
	UpdateFunc              func(ctx context.Context, token *domain.Token) error
	DeleteFunc              func(ctx context.Context, id uint) error
	DeleteByUserIDFunc      func(ctx context.Context, userID uint) error
	DeleteExpiredTokensFunc func(ctx context.Context) error
}

func (m *MockTokenRepository) Create(ctx context.Context, token *domain.Token) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, token)
	}
	return nil
}

func (m *MockTokenRepository) FindByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error) {
	if m.FindByAccessTokenFunc != nil {
		return m.FindByAccessTokenFunc(ctx, accessToken)
	}
	return nil, nil
}

func (m *MockTokenRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	if m.FindByRefreshTokenFunc != nil {
		return m.FindByRefreshTokenFunc(ctx, refreshToken)
	}
	return nil, nil
}

func (m *MockTokenRepository) Update(ctx context.Context, token *domain.Token) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, token)
	}
	return nil
}

func (m *MockTokenRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	if m.DeleteByUserIDFunc != nil {
		return m.DeleteByUserIDFunc(ctx, userID)
	}
	return nil
}

func (m *MockTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	if m.DeleteExpiredTokensFunc != nil {
		return m.DeleteExpiredTokensFunc(ctx)
	}
	return nil
}

// MockTokenService is a mock implementation of ports.TokenService for testing
type MockTokenService struct {
	GenerateTokenPairFunc    func(userID uint) (*domain.TokenPair, error)
	ValidateTokenFunc        func(token string) (uint, error)
	ValidateRefreshTokenFunc func(token string) (uint, error)
}

func (m *MockTokenService) GenerateTokenPair(userID uint) (*domain.TokenPair, error) {
	if m.GenerateTokenPairFunc != nil {
		return m.GenerateTokenPairFunc(userID)
	}
	return nil, nil
}

func (m *MockTokenService) ValidateToken(token string) (uint, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return 0, nil
}

func (m *MockTokenService) ValidateRefreshToken(token string) (uint, error) {
	if m.ValidateRefreshTokenFunc != nil {
		return m.ValidateRefreshTokenFunc(token)
	}
	return 0, nil
}

// MockAuthService is a mock implementation of ports.AuthService for testing
type MockAuthService struct {
	RegisterFunc            func(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	LoginFunc               func(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshTokenFunc        func(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error)
	LogoutFunc              func(ctx context.Context, accessToken string) error
	ValidateAccessTokenFunc func(ctx context.Context, accessToken string) (uint, error)
}

func (m *MockAuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthService) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthService) Logout(ctx context.Context, accessToken string) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, accessToken)
	}
	return nil
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, accessToken string) (uint, error) {
	if m.ValidateAccessTokenFunc != nil {
		return m.ValidateAccessTokenFunc(ctx, accessToken)
	}
	return 0, nil
}
